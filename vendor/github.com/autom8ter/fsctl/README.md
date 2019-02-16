![autom8ter](https://github.com/autom8ter/logo/blob/master/white_logo_dark_background.jpg?raw=true)

# FsCtl

**FsCtl is a package for managing an os filesystem**

**Author:** Coleman Word coleman.word@autom8ter.com

**Download**: `go get github.com/autom8ter/fsctl`

**License:** MIT

```text

 ______    _____ _   _ 
|  ____|  / ____| | | |
| |__ ___| |    | |_| |
|  __/ __| |    | __| |
| |  \__ \ |____| |_| |
|_|  |___/\_____|\__|_|
                       
                       
                                   
```

## Features
- embedded viper object
- embedded afero fs object
- synced environmental and configuration variables
- copy file
- stream file
- render templates from config (with sprigs excellent funcmap package) ref: https://github.com/Masterminds/sprig
- scan & replace

## ENV Variables:
CFGURL- remote github repot containing config file(s)
CFGToken- github token
GITUSER- github username
