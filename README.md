# space-envoy-service

```
go get github.com/niubirbang/space-envoy-service-go
```

```golang
import (
  service "github.com/niubirbang/space-envoy-service-go"
)
m := service.NewManager(serviceName, serviceFile)
m.Init()
m.Up(configDir, configFile)
m.Down()
m.Uninstall()
```
