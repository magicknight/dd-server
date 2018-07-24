package favorite

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/util"
	"github.com/shmy/dd-server/model/favorite"
	"github.com/shmy/dd-server/handler/middleware/jwt"
	"github.com/globalsign/mgo/bson"
	"time"
	"errors"
	"github.com/shmy/dd-server/model/video"
	"github.com/shmy/dd-server/model/collection"
)

// 获取所有收藏夹
func All (c echo.Context) error {
	cc := util.ApiContext{ c }
	user := cc.Get("user")
	userClaims := user.(*jwt.ClienJwtClaims)
	r, err := favorite.M.Find(bson.M{
		"_uid": userClaims.Id,
	}, nil)
	if err != nil {
		return cc.Fail(err)
	}
	return cc.Success(r)
}
// 创建一个收藏夹
func Create (c echo.Context) error {
	cc := util.ApiContext{ c }
	user := cc.Get("user")
	userClaims := user.(*jwt.ClienJwtClaims)
	name := cc.DefaultFormValueString("name", "", true)
	if name == "" {
		return cc.Fail(errors.New("请输入收藏夹名称"))
	}
	data := bson.M{
		"_id": bson.NewObjectId(),
		"name": name,
		"_uid": userClaims.Id,
		"created_at": time.Now(),
	}
	_, err := favorite.M.Insert(data)
	if err != nil {
		return cc.Fail(err)
	}
	return cc.Success(data)
}

// 添加一个资源到收藏夹

func AddToFavorite (c echo.Context) error {
	cc := util.ApiContext{ c }
	user := cc.Get("user")
	userClaims := user.(*jwt.ClienJwtClaims)
	_vid := cc.DefaultFormValueString("vid", "", true)
	_fid := cc.DefaultFormValueString("fid", "", true)
	if _vid == "" {
		return cc.Fail(errors.New("请输入视频id"))
	}
	if !bson.IsObjectIdHex(_vid) {
		return cc.Fail(errors.New("视频id格式不正确"))
	}
	if _fid == "" {
		return cc.Fail(errors.New("请输入收藏夹id"))
	}
	if !bson.IsObjectIdHex(_fid) {
		return cc.Fail(errors.New("收藏夹id格式不正确"))
	}
	// 判断视频是否存在
	v, err := video.M.FindById(_vid, "_id")
	if err != nil {
		return cc.Fail(err)
	}
	if v == nil {
		return cc.Fail(errors.New("视频不存在"))
	}
	// 判断收藏夹是否存在
	f, err := favorite.M.FindById(_fid, "_uid")
	if err != nil {
		return cc.Fail(err)
	}
	if f == nil {
		return cc.Fail(errors.New("收藏夹不存在"))
	}
	// 判断是否是该用户的收藏夹
	_uid := userClaims.Id
	if f["_uid"] != _uid {
		return cc.Fail(errors.New("收藏夹不属于你"))
	}
	data := bson.M{
		//"_id": bson.NewObjectId(),
		"_vid": _vid,
		"_fid": _fid,
		"_uid": _uid,
		//"created_at": time.Now(),
	}
	// 判断该人该收藏夹是否已经收藏过该视频了
	count, err := collection.M.Count(data)
	if err != nil {
		return cc.Fail(err)
	}
	if count != 0 {
		return cc.Fail(errors.New("该收藏夹已收藏过该视频了"))
	}
	data["_id"] = bson.NewObjectId()
	data["created_at"] = time.Now()
	_, err = collection.M.Insert(data)
	if err != nil {
		return cc.Fail(err)
	}
	return cc.Success(data)
}