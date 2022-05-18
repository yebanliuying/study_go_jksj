package week2

// ## 作业
// 我们在数据库操作的时候，比如 `dao` 层中当遇到一个 `sql.ErrNoRows` 的时候，是否应该 `Wrap` 这个 `error`，抛给上层。为什么？应该怎么做请写出代码

// ## 回答
// 应该。如果是一个整体项目处理error交给最顶级的调用层去记录，剩下的都去wrap一下错误，然后返回上一层。而且查询数据有可能是查询失败，可能是请求超时、错误查询等等或者是没有该数据。如果是三方库，就直接及时把错误抛出就行了。

// ## 伪代码：
type Dao struct{}

func (d *Dao) getData () (data *userModel, err error) {
	_, err = DB.Table(xx).select()
	if err != nil {
		return errors.Wrap(err, "getData failed")
	}
}

type Server struct{}

func (s *Server) getUser (data *userServer, err error) {
	data, err = dao.getData()
	if errors.Is(err, dao.ErrNoRows) {
		return nil
	}
}
