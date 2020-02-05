### 启动
- app.go 在init函数中注册配置文件启动器、数据库启动器、验证启动器、WebApi启动器、GoRPC启动器、RpcApi启动器、过期红包退回定期任务启动器、Hook启动器、Eureka启动器、Iris启动器。
- main.go 读取启动配置文件，从中获取consul中心地址，然后从consul中获取所需的共享配置信息。使用配置对象初始化一个BootApplication。调用BootApplication对象的Start方法开始逐步调用各启动器的Init、SetUp、Start方法。

### 基础设施
- boot.go 初始化一个BootApplication对象，该对象保存一个StarterContext，其中存有配置对象，用于多个启动器共享。该对象的就是用于管理各启动器初始化与启动的。
- initialier.go 定义一个初始化管理器接口，主要供上层Web Api注册使用，资金账户对象和红包对象把自身注册入该管理器，WebApi启动器会在SetUp阶段对该管理器内所有对象进行初始化（主要为请求Handler的注册）。
- starter.go 整个程序的启动核心，定义了一套程序启动初始化接口，各启动器实现其中一个或多个接口，其他接口默认处理。其中也有StarterContext用于各启动器共享数据。各启动器通过Register方法把自身注册进StarterRegister中。在boot.go中就会取出StarterRegister总保存的启动器来初始化。
- web.go 一个启动器，用于调用上层Web Api的Init初始化函数，进行Handler的注册。


#### base
- base.go 检查对象是否已经实例化
- dbx.go 使用配置实例化dbx数据库对象。该数据库具备返回隐藏对象的方法。
- dbx_base.go 该源码主要用于转调dbx的事务方法;把事务runner通过context传递，实现共用同一runner的事务要么全完成，要么全撤销。
- eureka.go 使用配置实例化eureka客户端对象，其中在InstanceInfo的MetaData中设置附加信息：rpc端口。然后该客户端向eureka服务器注册自身信息。
- hook.go 用于注册各启动器Stop方法，在程序接收到信号的时候通过回调这些方法正确优雅地关闭程序。该源码会注册需要处理的信号，然后通过从Channel中取得信号来进行关闭前的处理。
- iris_server.go 初始化一个iris程序对象。该对象使用了iris的Logger中间件，用于对到来的请求进行信息展示;使用了iris的recover中间件，用于处理请求时遇到panic错误的时候也不至于退出程序。最后通过iris的run方法开始进行请求监听。
- log.go 进行log展示的设置。如展示格式、文字颜色、高亮。
- props.go 取出程序启动时保存在StarterContext中实例化的配置对象并保存作包内可见，从配置中获取系统账户信息。
- res.go 定义一个作为HTTP response信息的结构体，包括：状态码（自定义）、信息、数据。
- rpc.go 实例化一个rpc server对象，在其中开始监听指定端口，并接受来自远端的rpc请求。开放接口给上层注册RPC Api。
- validator.go 实例化一个拦截器对象，用于验证请求的结构体是否符合要求。

#### gorpc
- rpc.go 在本地服务注册表中找到符合要求的微服务应用，然后再通过该应用指定的负载均衡算法获取一个应用实例，获取该实例RPC地址与端口，然后使用RPC客户端进行RPC接口调用。

#### httpclient
- http_client.go 提供实例化一个HttpClient对象接口，通过该对象的NewRequest方法来进行HTTP请求的生成（在本地服务注册表中找到符合要求的微服务应用，然后再通过该应用指定的负载均衡算法获取一个应用实例），通过该对象的Do方法来进行请求的发送与响应的处理。

#### lb
- app.go 在本地服务注册表中找到符合要求的微服务应用，对获取到的应用进行实例提取，提取出来后放入切片;通过各应用的负载均衡算法在切片中找出应用的实例。
- lb.go 定义一个负载均衡接口。
- lb_hash.go 实现负载均衡接口-hash
- lb_rr.go 实现负载均衡接口-轮询、随机

### 资金账户子系统
#### service
- accounts.go 定义了资金账户接口提供给应用层调用，也定义了一些DTO用于与上层应用进行数据交互。API：创建账户、转账、充值、通过用户ID获取红包账户、通过帐号编号获取账户。
- accounts_consts.go 定义了服务层一些常量。

#### core
##### accounts
- dao_account.go 数据库交互。通过账户编号获取账户、通过用户ID获取账户、插入新账户、账户余额的更新、账户状态更新。其中余额更新是使用了乐观锁。
- dao_account_log.go 数据库交互。通过流水编号获取账户流水、通过交易编号获取账户流水、插入账户流水。
- po_account.go 对象持久化。定义了账户的结构体，与数据库表中account表字段一一对应。还有该结构体与业务层数据交互DTO之间的转换。
- po_account_log.go 同上。
- domain_account.go 这个主要用于把某些操作整合在一起提供给外一个Api，比如Create方法创建一个账户的同时会创建账户流水。这里会调用ToDTO和FromDTO方法来转换数据。应该在一个事务中完成的操作都放到Tx中。提供创建账户、转账、获取账户、获取流水接口。
- service.go 实例化一个服务层对象，该对象实现了服务层账户接口。这里是使用应用层传递下来的数据DTO，返回时也是返回DTO数据。这里是对数据进行处理后，简单调用Domain层接口就行。

#### apis
##### web
- account.go 这里是定义HTTP handler的地方，用于提供RESTful Api。主要就是从Request中取出数据，调用service接口，然后返回数据，数据是Res结构体。

### 红包子系统