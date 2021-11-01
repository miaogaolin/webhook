# webhook
github webhook 工具，指定自动执行的脚本。

# 安装
```shell
go get github.com/miaogaolin/webhook
```

# 配置

创建 hooks.json 配置文件。
```json
{
  "bind": ":9000",
  "items": [
    {
      "repo": "https://github.com/miaogaolin/printlove",
      "branch": "main",
      "script": "deploy.sh",
      "secret": "123123"
    }
  ]
}
```

# 启动
```shell
webhook hooks.json
```

