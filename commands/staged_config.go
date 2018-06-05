package commands

import (
	"fmt"

	"github.com/pivotal-cf/jhanda"
	"github.com/pivotal-cf/om/api"
	"gopkg.in/yaml.v2"
	"strings"
	"strconv"
)

type StagedConfig struct {
	logger  logger
	service stagedConfigService
	Options struct {
		Product            string `long:"product-name" short:"p" required:"true" description:"name of product"`
		IncludeCredentials bool   `short:"c" long:"include-credentials" description:"include credentials. note: requires product to have been deployed"`
		IncludePlaceholder bool   `short:"r" long:"include-placeholder" description:"replace obscured credentials to interpolatable placeholder"`
	}
}

//go:generate counterfeiter -o ./fakes/staged_config_service.go --fake-name StagedConfigService . stagedConfigService
type stagedConfigService interface {
	GetDeployedProductCredential(input api.GetDeployedProductCredentialInput) (api.GetDeployedProductCredentialOutput, error)
	GetStagedProductByName(product string) (api.StagedProductsFindOutput, error)
	GetStagedProductJobResourceConfig(productGUID, jobGUID string) (api.JobProperties, error)
	GetStagedProductNetworksAndAZs(product string) (map[string]interface{}, error)
	GetStagedProductProperties(product string) (map[string]api.ResponseProperty, error)
	ListDeployedProducts() ([]api.DeployedProductOutput, error)
	ListStagedProductJobs(productGUID string) (map[string]string, error)
}

func NewStagedConfig(service stagedConfigService, logger logger) StagedConfig {
	return StagedConfig{
		logger:  logger,
		service: service,
	}
}

func (ec StagedConfig) Usage() jhanda.Usage {
	return jhanda.Usage{
		Description:      "This command generates a config from a staged product that can be passed in to om configure-product (Note: credentials are not available and will appear as '***')",
		ShortDescription: "**EXPERIMENTAL** generates a config from a staged product",
		Flags:            ec.Options,
	}
}

func (ec StagedConfig) Execute(args []string) error {
	if _, err := jhanda.Parse(&ec.Options, args); err != nil {
		return fmt.Errorf("could not parse staged-config flags: %s", err)
	}

	if ec.Options.IncludeCredentials {
		deployedProducts, err := ec.service.ListDeployedProducts()
		if err != nil {
			return err
		}
		var productDeployed bool
		for _, p := range deployedProducts {
			if p.Type == ec.Options.Product {
				productDeployed = true
				break
			}
		}
		if !productDeployed {
			return fmt.Errorf("cannot retrieve credentials for product '%s': deploy the product and retry", ec.Options.Product)
		}
	}

	findOutput, err := ec.service.GetStagedProductByName(ec.Options.Product)
	if err != nil {
		return err
	}
	productGUID := findOutput.Product.GUID

	properties, err := ec.service.GetStagedProductProperties(productGUID)
	if err != nil {
		return err
	}

	configurableProperties := map[string]interface{}{}
	selectorProperties := map[string]string{}

	for name, property := range properties {
		if property.Configurable && property.Value != nil {
			if property.Type == "selector" {
				selectorProperties[name] = property.Value.(string)
			}
			output, err := ec.parseProperties(productGUID, name, property)
			if err != nil {
				return err
			}
			if output != nil && len(output) > 0 {
				configurableProperties[name] = output
			}
		}
	}

	for name := range configurableProperties {
		components := strings.Split(name, ".")[1:] // the 0th item is an empty string due to `.some.other`
		if len(components) == 2 {
			continue
		}
		selector := "." + strings.Join(components[:2], ".")
		if val, ok := selectorProperties[selector]; ok && components[2] != val {
			delete(configurableProperties, name)
		}
	}

	networks, err := ec.service.GetStagedProductNetworksAndAZs(productGUID)
	if err != nil {
		return err
	}

	jobs, err := ec.service.ListStagedProductJobs(productGUID)
	if err != nil {
		return err
	}

	resourceConfig := map[string]api.JobProperties{}

	for name, jobGUID := range jobs {
		jobProperties, err := ec.service.GetStagedProductJobResourceConfig(productGUID, jobGUID)
		if err != nil {
			return err
		}

		resourceConfig[name] = jobProperties
	}

	config := struct {
		Properties               map[string]interface{}       `yaml:"product-properties"`
		NetworkProperties        map[string]interface{}       `yaml:"network-properties"`
		ResourceConfigProperties map[string]api.JobProperties `yaml:"resource-config"`
	}{
		Properties:               configurableProperties,
		NetworkProperties:        networks,
		ResourceConfigProperties: resourceConfig,
	}

	output, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %s", err) // un-tested
	}
	ec.logger.Println(string(output))

	return nil
}

func (ec StagedConfig) parseProperties(productGUID string, name string, property api.ResponseProperty) (map[string]interface{}, error) {
	if property.IsCredential {
		return ec.handleCredential(productGUID, name, property)
	} else if property.Type == "collection" {
		var valueItems []map[string]interface{}
		for index, valueItem := range property.Value.([]interface{}) {
			valueItemTyped := valueItem.(map[interface{}]interface{})
			tempMap := make(map[string]interface{})
			for itemKey, itemVal := range valueItemTyped {
				itemValTyped := itemVal.(map[interface{}]interface{})

				apiRes := api.ResponseProperty{
					Value: itemValTyped["value"],
					Configurable: itemValTyped["configurable"].(bool),
					IsCredential: itemValTyped["credential"].(bool),
					Type: itemValTyped["type"].(string),
				}
				valueNamePrefix := name + "[" + strconv.Itoa(index) + "]." + itemKey.(string)
				retVal, err := ec.parseProperties(productGUID, valueNamePrefix, apiRes)
				if err != nil {
					return nil, err
				}
				if retVal != nil && len(retVal) > 0 {
					tempMap[itemKey.(string)] = retVal
				}
			}
			if len(tempMap) >0 {
				valueItems = append(valueItems, tempMap)
			}
		}
		if len(valueItems) > 0 {
			return map[string]interface{}{"value": valueItems}, nil
		} else {
			return nil, nil
		}
	} else {
		return map[string]interface{}{"value": property.Value}, nil
	}
	return nil, nil
}

func (ec StagedConfig) handleCredential(productGUID string, name string, property api.ResponseProperty) (map[string]interface{}, error) {
	var output map[string]interface{}

	if ec.Options.IncludeCredentials {
		apiOutput, err := ec.service.GetDeployedProductCredential(api.GetDeployedProductCredentialInput{
			DeployedGUID:        productGUID,
			CredentialReference: name,
		})
		if err != nil {
			return nil, err
		}
		output = map[string]interface{}{"value": apiOutput.Credential.Value}
	} else if ec.Options.IncludePlaceholder {
		switch property.Type {
		case "secret":
			output = map[string]interface{}{
				"value": map[string]string{
					"secret": fmt.Sprintf("((%s.secret))", name),
				},
			}
		case "simple_credentials":
			output = map[string]interface{}{
				"value": map[string]string{
					"identity": fmt.Sprintf("((%s.identity))", name),
					"password": fmt.Sprintf("((%s.password))", name),
				},
			}
		case "rsa_cert_credentials":
			output = map[string]interface{}{
				"value": map[string]string{
					"cert_pem":        fmt.Sprintf("((%s.cert_pem))", name),
					"private_key_pem": fmt.Sprintf("((%s.private_key_pem))", name),
				},
			}
		case "rsa_pkey_credentials":
			output = map[string]interface{}{
				"value": map[string]string{
					"private_key_pem": fmt.Sprintf("((%s.private_key_pem))", name),
				},
			}
		case "salted_credentials":
			output = map[string]interface{}{
				"value": map[string]string{
					"identity": fmt.Sprintf("((%s.identity))", name),
					"password": fmt.Sprintf("((%s.password))", name),
					"salt":     fmt.Sprintf("((%s.salt))", name),
				},
			}
		}
	} else {
		output = nil
	}

	return output, nil
}

func addSecretPlaceholder(value interface{}, t string, configurableProperties map[string]interface{}, name string) {
	switch t {
	case "secret":
		configurableProperties[name] = map[string]interface{}{
			"value": map[string]string{
				"secret": fmt.Sprintf("((%s.secret))", name),
			},
		}
	case "simple_credentials":
		configurableProperties[name] = map[string]interface{}{
			"value": map[string]string{
				"identity": fmt.Sprintf("((%s.identity))", name),
				"password": fmt.Sprintf("((%s.password))", name),
			},
		}
	case "rsa_cert_credentials":
		configurableProperties[name] = map[string]interface{}{
			"value": map[string]string{
				"cert_pem":        fmt.Sprintf("((%s.cert_pem))", name),
				"private_key_pem": fmt.Sprintf("((%s.private_key_pem))", name),
			},
		}
	case "rsa_pkey_credentials":
		configurableProperties[name] = map[string]interface{}{
			"value": map[string]string{
				"private_key_pem": fmt.Sprintf("((%s.private_key_pem))", name),
			},
		}
	case "salted_credentials":
		configurableProperties[name] = map[string]interface{}{
			"value": map[string]string{
				"identity": fmt.Sprintf("((%s.identity))", name),
				"password": fmt.Sprintf("((%s.password))", name),
				"salt":     fmt.Sprintf("((%s.salt))", name),
			},
		}
	case "collection":
		collectionValue := value.([]interface{})
		finalValue := []interface{}{}
		for index, item := range collectionValue {
			collectionValueItems := item.(map[string]api.ResponseProperty)

			for key, val := range collectionValueItems {
				collectionName := name + "." + strconv.Itoa(index) + "." + key
				placeholder := map[string]interface{}{}
				addSecretPlaceholder(val.Value, val.Type, placeholder, collectionName)

				for _, temp2 := range placeholder {
					temp2Typed := temp2.(map[string]interface{})
					finalValue = append(finalValue, map[string]interface{}{
						key: map[string]interface{}{
							"value": temp2Typed["value"],
						},
					})
				}

			}
			configurableProperties[name] = map[string]interface{}{
				"value": finalValue,
			}
		}
	default:
		configurableProperties[name] = map[string]interface{}{"value": value}
	}
}
