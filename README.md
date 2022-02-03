# Love Fortune

Love Fortune 是一个定时执行 [Bark](https://github.com/Finb/Bark) 任务的应用。他提供了一系列 REST API 接口可供用户添加、删除 Bark 任务。同时 Love Fortune 支持模板消息，由 `plugin` 模块提供类似于纪念日、每日天气之类的模板数据。

## Run

运行前需要设置一个`ACCESS_KEY`用于保护 API 访问端点。

```sh
export ACCESS_KEY=40be1d22-9924-4409-b86e-2d095b22fece
./lovefortune
```

运行后会在当前目录下生成数据库文件：`lovefortune.db`所有的定时任务都会保存在此数据库中。

## API Example

添加任务

> `spec`字段的格式参考[CRON Expression Format](https://pkg.go.dev/github.com/robfig/cron?utm_source=godoc#hdr-CRON_Expression_Format)。

```sh
curl --request POST \
  --url http://localhost:8080/api/tasks \
  --header 'X-Access-Token: 40be1d22-9924-4409-b86e-2d095b22fece' \
  --header 'content-type: application/json' \
  --data '{
  "title": "纪念日",
  "content": "今天是我们在一起第{{.LoveAnniversaryDays}}天啦！",
  "icon": null,
  "deviceKey": "WSsV6scqP2VYrTQ9M46YkF",
  "description": "纪念日(每天执行)",
  "spec": "0 22 21 * * *"
}'
```

列出所有任务

```sh
curl --request GET \
  --url http://localhost:8080/api/tasks \
  --header 'X-Access-Token: 40be1d22-9924-4409-b86e-2d095b22fece'
```

删除任务

```sh
curl --request DELETE \
  --url http://localhost:8080/api/tasks/3467d497-4629-4926-a233-4aaf804c3af9 \
  --header 'X-Access-Token: 40be1d22-9924-4409-b86e-2d095b22fece'
```

查看某任务的所有运行日志

```sh
curl --request GET \
  --url http://localhost:8080/api/logs/1e0adbf1-4b52-4140-a6cc-9c50d3e5e669 \
  --header 'X-Access-Token: 40be1d22-9924-4409-b86e-2d095b22fece' \
  --header 'content-type: application/json'
```
