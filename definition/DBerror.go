// 对数据库操作的响应码
package definition

type DBcode string

const (
	DB_ERROR                 DBcode = "100" // 操作失败,一般是未知错误
	DB_ERROR_PARAM                  = "101" // 操作失败,入参不合法
	DB_ERROR_UNEXPECTED             = "102" // 操作失败,下层抛出了本层捕捉不到的错误，可能是代码有变更
	DB_ERROR_DATA_FMT               = "103" // 数据的存储格式有误，解析出错
	DB_ERROR_TX                     = "104" // 事务执行失败,具体看日志 err
	DB_ERROR_UNAME_UNIQUE           = "110" // 操作失败,用户名唯一
	DB_ERROR_NICKNAME_UNIQUE        = "111" // 操作失败,昵称唯一

	DB_SUCCESS           DBcode = "200" // 操作成功
	DB_EXIST             DBcode = "201" // 要查询的数据存在
	DB_NOEXIST           DBcode = "210" // 要查询的数据不存在
	DB_NOEXIST_USER      DBcode = "211" // 用户不存在
	DB_NOEXIST_POST      DBcode = "212" // 帖子不存在
	DB_NOEXIST_TARGET    DBcode = "213" // 回复目标（评论或回复私聊对象）不存在
	DB_NOEXIST_ADDRESSEE DBcode = "214" // 回复目标（私聊对象）不存在

	DB_UNMATCH       DBcode = "400" // id对不上,没有操作权限
	DB_NOT_THE_OWNER DBcode = "401" // 隐私信息，不是本人不允许查看

)
