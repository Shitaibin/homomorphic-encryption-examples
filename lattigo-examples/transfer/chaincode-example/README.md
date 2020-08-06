基于同态加密的转账Demo
============

分为3个项目:
- chaincode目录: 链码项目
- api-server: 提供通过RESTful API调用链码的接口服务器
- front-end: 前端页面


## API-Server功能

1. 创建银行 [✔️]
    1. 创建银行请求，Server创建银行同态加密密钥，公钥和私钥，保存到内存（长期需要持久化到文件或DB）。请求参数为BankID
    2. 调用链码`AddBankPublicKey`把银行的公钥上链，以便对用户数据进行同态加密的修改
    3. 请求响应为：BankID，SK、PK太长，不返回
2. 设置银行用户账户余额 [✔]
    1. 请求参数为：BankID、AccountID、Balance（明文）
    2. Server利用该银行的公钥对Amount进行加密得到余额密文CipherBalance
    3. Server调用`SetAccountBalance`把CipherBalance上链
    4. 响应为：BankID、AccountID、Status（成功或失败）
3. 查询银行用户账户余额 [✔]
    1. 请求参数为：BankID、AccountID
    2. Server调用链码`QueryAccountBalance`，获取用户余额，结果为同态加密的用户余额CipherBalance
    3. Server利用该银行的私钥对CipherBalance进行解密，获得Balance
    4. 响应为：BankID、AccountID、Balance
4. 用户余额链上转账 [✔️]
    1. 请求参数为：FromBankID、FromAccountID、ToBankID、ToAccountID、Amount
    2. Server调用链码`Transfer`进行链上用户余额转账
    3. 响应为：FromBankID、FromAccountID、ToBankID、ToAccountID、Amount、Status（成功或失败）
5. 银行密钥文件下载 [✔️]

## 演示交互案例

部署：
1. 在区块链上部署好转账链码
2. 启动2个API-Server，每个Server供1个银行使用，Server1连接区块链Org1的节点，Server2连接Org2的节点

演示操作流程：
1. 向Server1发送创建银行1的请求
2. 向Server1发送设置账户1余额的请求，余额100
3. 向Server1发送查询账户1余额的请求，校对余额是否为100
3. 向Server2发送创建银行2的请求
4. 向Server2发送设置账户2余额的请求，余额100
3. 向Server2发送查询账户2余额的请求，校对余额是否为100
1. 向Server1发送转账请求，实现从账户1向银行2账户2的转账，金额30
1. 向Server1发送查询账户1余额的请求，校对余额是否为70
1. 向Server2发送查询账户2余额的请求，校对余额是否为130