# nftserver
因提供的nft库只有java、linux版本的，并没有hp-ux版本下的。为此，需要搭建一个服务来接收从hp上转来的nft命令请求并转发给
nft服务。这里提供了一个golang版本下的实现。主要实现细节:
1、通过监听服务的方式，接收远程服务发来的nft命令请求；
2、将命令解析后，通过cgo调用nft提供的linux库，实现nft的相关任务，如任务增加、查询、取消。
本版本作为benchmark的参照物，用来评估其他实现方式的性能。
