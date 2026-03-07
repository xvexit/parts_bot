package part

type Part struct{
	ID int
	Name string
	Price int
	DeliveryDay int
	ImgUrl string
}

func NewPart(
	id, price, delivDay int,
	name, imgUrl string,
	) *Part{
	return &Part{
		ID: id,
		Name: name,
		Price: price,
		DeliveryDay: delivDay,
		ImgUrl: imgUrl,
	}
}