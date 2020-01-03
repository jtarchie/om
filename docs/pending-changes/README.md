<!--- This file is autogenerated from the files in docsgenerator/templates/pending-changes --->
&larr; [back to Commands](../README.md)

# `om pending-changes`

This authenticated command lists all products and will display whether they are unchanged (no pending changes) or changed (has pending changes).

## Command Usage
```
ॐ  pending-changes
This authenticated command lists all products and will display whether they are unchanged (no pending changes) or changed (has pending changes).

Usage: om [options] pending-changes [<args>]
  --ca-cert, OM_CA_CERT                                  string  OpsManager CA certificate path or value
  --client-id, -c, OM_CLIENT_ID                          string  Client ID for the Ops Manager VM (not required for unauthenticated commands)
  --client-secret, -s, OM_CLIENT_SECRET                  string  Client Secret for the Ops Manager VM (not required for unauthenticated commands)
  --connect-timeout, -o, OM_CONNECT_TIMEOUT              int     timeout in seconds to make TCP connections (default: 10)
  --decryption-passphrase, -d, OM_DECRYPTION_PASSPHRASE  string  Passphrase to decrypt the installation if the Ops Manager VM has been rebooted (optional for most commands)
  --env, -e                                              string  env file with login credentials
  --help, -h                                             bool    prints this usage information (default: false)
  --password, -p, OM_PASSWORD                            string  admin password for the Ops Manager VM (not required for unauthenticated commands)
  --request-timeout, -r, OM_REQUEST_TIMEOUT              int     timeout in seconds for HTTP requests to Ops Manager (default: 1800)
  --skip-ssl-validation, -k, OM_SKIP_SSL_VALIDATION      bool    skip ssl certificate validation during http requests (default: false)
  --target, -t, OM_TARGET                                string  location of the Ops Manager VM
  --trace, -tr, OM_TRACE                                 bool    prints HTTP requests and response payloads
  --username, -u, OM_USERNAME                            string  admin username for the Ops Manager VM (not required for unauthenticated commands)
  --version, -v                                          bool    prints the om release version (default: false)
  OM_VARS_ENV                                            string  **EXPERIMENTAL** load vars from environment variables by specifying a prefix (e.g.: 'MY' to load MY_var=value)

Command Arguments:
  --check       bool    Exit 1 if there are any pending changes. Useful for validating that Ops Manager is in a clean state.
  --format, -f  string  Format to print as (options: table,json) (default: table)

```
