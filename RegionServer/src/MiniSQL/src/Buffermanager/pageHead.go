package Buffermanager

type PageHead struct {
	pageID  int32 //页的编号
	isFixed bool  //是否常驻内存
}
