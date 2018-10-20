# devmp(Device management platform)
> 设备管理平台

- 用于管理公司办公电脑
- 查询办公电脑相关信息
- 年综汇总数据导出编辑
- 学习如何使用sqlite3数据库
- 导出所有设备信息并且将每个设备信息转换成二维码方便以后贴在设备上

为什么把网页模板拆分的那么零散的原因就是想看看模板嵌套的方法,以及如何向子模板里面传递数据.


## 时序图

> 查询数据
```sequence
操作者-服务器: 给你MAC,我要信息.
Note right of 操作者: get方法发送mac地址
服务器->服务器: 检查数据.
Note left of 服务器: 数据格式不正确
服务器-->操作者: 您输入的格式不正确
Note right of 服务器: 数据格式正确
服务器->数据库: 快点帮我查查这东西在不在.
数据库->数据库: 查询数据
Note left of 数据库: 查询失败
数据库-->服务器: 没有对应数据.
服务器-->操作者: 查询失败.
Note right of 数据库: 查询成功
数据库-->服务器: 你要的东西在这里,拿去吧.
服务器-->操作者: 好的给你你要的数据.
```
> 添加、修改数据
```sequence
操作者->服务器: 我突然想要添加一个设备 .
服务器-->操作者: 新的页面打开,请填写内容 .
操作者->服务器: 这是我的数据 .
服务器->服务器: 检查数据.
Note left of 服务器: 数据格式不正确.
服务器-->操作者: 你输入的格式不正确
note left of 操作者: 重新输入并且提交.
Note right of 服务器: 数据格式正确.
服务器->数据库: 请添加这个新的设备信息.
数据库-数据库: 查询数据中...
Note left of 数据库: 已有设备
数据库-->服务器: 已存在该设备,修改信息吗
服务器-->操作者: 覆盖之前的设备信息吗
操作者-操作者: 决定是否覆盖
note left of 操作者: 放弃则退出.
操作者->服务器: 是的覆盖
服务器->数据库: 覆盖
数据库-数据库: 写入中
数据库-->服务器: 写入成功.
服务器-->操作者: 写入成功
Note right of 数据库: 没有设备
数据库-数据库: 写入中
数据库-->服务器: 写入成功
服务器-->操作者: 添加完成,以下是你刚才添加的数据请确认


```

## 流程图:

```flow
st=>start: 打开首页
e=>end
s=>condition: 是否查询
imac=>inputoutput: 输入MAC
iedit=>inputoutput: 输入设备信息
edit=>operation: 点击添加按钮
sub1=>subroutine: 查询
q=>condition: 查询结果
or=>operation: 返回数据
of=>operation: 提示错误
sub2=>subroutine: 添加数据
t=>condition: 添加结果

st->s
s(yes)->imac
s(no)->edit
imac->sub1
sub1->q
q(yes)->or
q(no)->of
or->e
edit->iedit
iedit->sub2
sub2->t
t(yes)->or
t(no)->of
of->e
```


开始编写日期: 2018年10月14日
