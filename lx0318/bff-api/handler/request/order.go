package request

type CreateOrder struct {
	MemberId  int         `form:"memberId" json:"memberId" xml:"memberId"  binding:"required"`
	AddressId int         `form:"addressId" json:"addressId" xml:"addressId" binding:"required"`
	OrderItem []OrderItem `form:"orderItem" json:"orderItem" xml:"orderItem" binding:"required"`
}

type OrderItem struct {
	ProductId int `form:"productId" json:"productId" xml:"productId" binding:"required"`
	Quantity  int `form:"quantity" json:"quantity" xml:"quantity"  binding:"required"`
}
