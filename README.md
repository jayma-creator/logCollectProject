#背景 
每个系统都有日志，当系统出现问题时，需要通过日志解决问题
  当系统机器比较少时，登陆到服务器上查看即可满足
  当系统机器规模巨大，登陆到机器上查看几乎不现实
  每个业务系统都有自己的⽇志，当业务系统出现问题时，需要通过查找⽇志信息来定位和解决问题。 当业务系统服务器⽐较少时，登陆到服务器上查看即可满⾜。但当系统机器规模巨⼤，登陆到服务器上查看⼏乎不现实（分布式的系统，⼀个系统部署在⼗⼏甚至几十台服务器上）

平常我们在进行业务开发时常常不免遇到下面几个问题:

当系统出现问题后，如何根据日志迅速的定位问题出在一个应用层？
在平常的工作中如何根据日志分析出一个请求到系统主要在那个应用层耗时较大？
在平常的工作中如何获取一个请求到达系统后在各个层测日志汇总？
针对以上问题，我们想要实现的一个解决方案是：

把机器上的日志实时收集，统一的存储到中心系统
然后再对这些日志建立索引，通过搜索即可以找到对应日志
通过提供界面友好的web界面，通过web即可以完成日志搜索
关于实现这个系统时可能会面临的问题：

实时日志量非常大，每天几十亿条
日志准实时收集，延迟控制在分钟级别
能够水平可扩展

##业界方案

有早期的ELK到现在的EFK。ELK在每台服务器上部署logstash，比较重量级，所以演化成客户端部署filebeat的EFK，由filebeat收集向logstash中写数据，最后落地到elasticsearch，通过kibana界面进行日志检索。


优点：现成的解决方案，直接拿过来用，能够实现日志收集与检索。

缺点：

运维成本⾼，每增加⼀个⽇志收集项，都需要⼿动修改配置
监控缺失，⽆法准确获取logstash的状态。⽆法做到定制化开发与维护
⽆法做到定制化开发与维护



- Log Agent，日志收集客户端，用来收集服务器上的日志,发往kafka
- Log Transfer从kafka中消费数据，然后写到es中
