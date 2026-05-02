package x11

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/composite"
	"github.com/jezek/xgb/xproto"
)

type Connection interface {
	Close()
	DefaultScreen() *xproto.ScreenInfo
	InternAtom(OnlyIfExists bool, NameLen uint16, Name string) InternAtomCookie
	GetProperty(Delete bool, Window xproto.Window, Property, Type xproto.Atom, LongOffset, LongLength uint32) GetPropertyCookie
	GetGeometry(Drawable xproto.Drawable) GetGeometryCookie
	TranslateCoordinates(SrcWindow, DstWindow xproto.Window, SrcX, SrcY int16) TranslateCoordinatesCookie
	SendEventChecked(Propagate bool, Destination xproto.Window, EventMask uint32, Event string) SendEventCookie
	MapWindowChecked(Window xproto.Window) MapWindowCookie
	ConfigureWindowChecked(Window xproto.Window, ValueMask uint16, ValueList []uint32) ConfigureWindowCookie
	NewId() (uint32, error)
	GetImage(Format byte, Drawable xproto.Drawable, X, Y int16, Width, Height uint16, PlaneMask uint32) GetImageCookie
	FreePixmap(Pixmap xproto.Pixmap) FreePixmapCookie
	InitComposite() error
	NameWindowPixmap(Window xproto.Window, Pixmap xproto.Pixmap) NameWindowPixmapCookie
}

type XgbConnection struct {
	conn *xgb.Conn
}

func NewXgbConnection(displayName string) (*XgbConnection, error) {
	conn, err := xgb.NewConnDisplay(displayName)
	if err != nil {
		return nil, err
	}
	return &XgbConnection{conn: conn}, nil
}

func (c *XgbConnection) Close() {
	c.conn.Close()
}

func (c *XgbConnection) DefaultScreen() *xproto.ScreenInfo {
	return xproto.Setup(c.conn).DefaultScreen(c.conn)
}

func (c *XgbConnection) InternAtom(OnlyIfExists bool, NameLen uint16, Name string) InternAtomCookie {
	cookie := xproto.InternAtom(c.conn, OnlyIfExists, NameLen, Name)
	return NewXProtoInternAtomCookie(cookie)
}

func (c *XgbConnection) GetProperty(Delete bool, Window xproto.Window, Property, Type xproto.Atom, LongOffset, LongLength uint32) GetPropertyCookie {
	cookie := xproto.GetProperty(c.conn, Delete, Window, Property, Type, LongOffset, LongLength)
	return NewXProtoGetPropertyCookie(cookie)
}

func (c *XgbConnection) GetGeometry(Drawable xproto.Drawable) GetGeometryCookie {
	cookie := xproto.GetGeometry(c.conn, Drawable)
	return NewXProtoGetGeometryCookie(cookie)
}

func (c *XgbConnection) TranslateCoordinates(SrcWindow, DstWindow xproto.Window, SrcX, SrcY int16) TranslateCoordinatesCookie {
	cookie := xproto.TranslateCoordinates(c.conn, SrcWindow, DstWindow, SrcX, SrcY)
	return NewXProtoTranslateCoordinatesCookie(cookie)
}

func (c *XgbConnection) SendEventChecked(Propagate bool, Destination xproto.Window, EventMask uint32, Event string) SendEventCookie {
	cookie := xproto.SendEventChecked(c.conn, Propagate, Destination, EventMask, Event)
	return NewXProtoSendEventCookie(cookie)
}

func (c *XgbConnection) MapWindowChecked(Window xproto.Window) MapWindowCookie {
	cookie := xproto.MapWindowChecked(c.conn, Window)
	return NewXProtoMapWindowCookie(cookie)
}

func (c *XgbConnection) ConfigureWindowChecked(Window xproto.Window, ValueMask uint16, ValueList []uint32) ConfigureWindowCookie {
	cookie := xproto.ConfigureWindowChecked(c.conn, Window, ValueMask, ValueList)
	return NewXProtoConfigureWindowCookie(cookie)
}

func (c *XgbConnection) NewId() (uint32, error) {
	return c.conn.NewId()
}

func (c *XgbConnection) GetImage(Format byte, Drawable xproto.Drawable, X, Y int16, Width, Height uint16, PlaneMask uint32) GetImageCookie {
	cookie := xproto.GetImage(c.conn, Format, Drawable, X, Y, Width, Height, PlaneMask)
	return NewXProtoGetImageCookie(cookie)
}

func (c *XgbConnection) FreePixmap(Pixmap xproto.Pixmap) FreePixmapCookie {
	cookie := xproto.FreePixmap(c.conn, Pixmap)
	return NewXProtoFreePixmapCookie(cookie)
}

func (c *XgbConnection) InitComposite() error {
	return composite.Init(c.conn)
}

func (c *XgbConnection) NameWindowPixmap(Window xproto.Window, Pixmap xproto.Pixmap) NameWindowPixmapCookie {
	cookie := composite.NameWindowPixmap(c.conn, Window, Pixmap)
	return NewXProtoNameWindowPixmapCookie(cookie)
}
