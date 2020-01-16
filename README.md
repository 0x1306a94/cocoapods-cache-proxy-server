# cocoapods-cache-proxy-server

#### CocoaPods 缓存服务 [CocoaPods 插件](https://github.com/0x1306a94/cocoapods-cache-proxy)

#### Homebrew
```shell
# install
brew tap 0x1306a94/homebrew-tap
brew install cocoapods-cache-proxy

# start
brew services start cocoapods-cache-proxy

# stop
brew services stop cocoapods-cache-proxy

# restart
brew services restart cocoapods-cache-proxy

# 修改配置文件
/usr/local/Cellar/cocoapods-cache-proxy-server/cache/cocoapods-cache-proxy-server/conf.yaml
```

#### Docker 运行
```shell
docker pull 0x1306a94/cocoapods-cache-proxy:v2
docker run -it -p 9898:9898  0x1306a94/cocoapods-cache-proxy:v2
```


