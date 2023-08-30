# USDT自动兑换TRX机器人

本程序由golang编写对USDT兑换TRX全天24自动处理，订单成功或失败将通知到群组

该程序运行于telegram机器人，需要准备机器人token。一键部署启动机器人

#### 功能
* 🔥 自定义利润率兑换比例
* 🔥 自定义余额提醒
* 🔥 支持预支付，自定义预支付金额
* 🔥 自动导入钱包私钥
* 🔥 不依赖第三方程序，安全可靠
* 🔥 成功通知，以及群组通知订单
* 🔥 订单失败支持管理员补发

#### 准备环境
确保机器人部署服务可以访问外网telegram.org

### 更多程序
* [telegram-trx](https://github.com/zavierswang/telegram-trx) **TRX兑换机器人**
* [telegram-monitor](https://github.com/zavierswang/telegram-monitor) **TRC20钱包事件监听机器人**
* [telegram-search](https://github.com/zavierswang/telegram-search) **导航机器人**（可支持全网搜索，API收费有点小贵）
* [telegram-premium](https://github.com/zavierswang/telegram-premium) **Telegram Premium自动充值机器人**
* [telegram-replay](https://github.com/zavierswang/telegram-replay) **双向机器人**
* [telegram-energy](https://github.com/zavierswang/telegram-energy) **TRON能量租凭机器人**
* [telegram-proto](https://github.com/zavierswang/telegram-proto) **Telegram协议号机器人**


### 部署
* 本程序基于`Telegram Bot`，主程序`telegram-trx`
* 确保机器人部署服务可以访问外网`telegram.org`
* 使用自己的`telegram`生成一个机器人，并获取到`token`
* 配置文件`telegram-trx.yaml.example`改名为`telegram-trx.yaml`, 修改建议配置项



> **注意：**
> * 不支持交易所转帐事件监听
> * 对linux不熟悉的给点打赏手把手教学🤭
> * 配置文件中的`license`配置请找 [🫣我](https://t.me/tg_llama) 拿~

