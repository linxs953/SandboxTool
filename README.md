# sandbox tool
 
 

## Install


```shell
git clone git@git.forchange.cn:linqinbo/sandboxtool.git
cd sandboxtool
go build  -o rtctl .
chmod u+x rtctl
cp rtctl /usr/local/bin
```



`支持的命令`

```shell
rtctl debug  [envirment] [unionid]
rtctl start [unionid] --it=true -v=1
rtctl start [puid] [storage_id] --puid=true --it=true
rtctl recycle [puid] --puid=true
rtctl recycle  --ep={sandbox_endpoint} --token={sandbox_token}
rtctl restart [unionid] --it=true -v=1
rtctl restart [puid] --it=true --puid=true
rtctl sdb --url='wss_url' --alive=true
rtctl sdb --ep="sanbox_endpoint" -t="sandbox_token"
```

参数含义
- start：  启动一个沙盒
    - `version`参数表示采用哪种方式获取沙盒，比如 1 表示用这个地址`entry/sandboxes/{sandboxid}`请求沙盒，2 表示用这个地址`entry/sandbox/{pool_name}`请求沙盒
    - `it`参数表示启动沙盒之后是否创建ws与沙盒通讯
- recycle：回收一个沙盒
- restart：强制重启一个沙盒
    - `version`语义与start相同，-v=1表示用这个`entry/sandboxes/{sandboxid}?force=1`重启，-v=2表示用这个``entry/sandbox/{pool_name}?force=1重启
- debug：  debug一个用户，根据[environment]跳转到指定的环境，[environment]可以为test或者prod
- sdb：    沙盒调试桥，支持`ep`和`url`两种方式, `alive`表示是否不断开连接，值为true或者false





