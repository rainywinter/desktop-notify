# http 桌面通知服务

在任务作业结束时发送一个桌面通知可以避免不停地切换窗口查看任务是否结束。尤其对一些下载任务、编译任务实用。

Usage:  

    make 启动服务


curl 发送通知

    curl localhost:8080/?message=任务完成