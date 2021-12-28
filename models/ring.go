package models

type Ring struct {
	list []interface{}
	index int
	maxLen int
	isFull bool
}

func InitRing(maxLen int)*Ring{
	return &Ring{
		list:   make([]interface{}, 0),
		index:  0,
		maxLen: maxLen,
		isFull: false,
	}
}

//向数据环插入数据
func (r *Ring)Insert(data interface{}){
	if !r.isFull{
		r.list = append(r.list, data)
		if len(r.list) == r.maxLen{
			r.isFull = true
		}
		return
	}

	r.list[r.index] = data
	if r.index >= r.maxLen-1{
		r.index = 0
	}else {
		r.index++
	}
}

//顺序返回数据列表；
//当请求数据量大于环中数据总量时，返回环最大数据量
func (r *Ring)GetList(num int)(list []interface{}){
	if !r.isFull{
		if num > len(r.list){
			num = len(r.list)
		}

	    list = r.list[:num]
		return
	}

	if r.maxLen - r.index >= num{
		list = r.list[r.index:r.index+num]
		return
	}

	if num >= r.maxLen{
		list = r.list[r.index:]
		copy(list, r.list[:r.index])
		return
	}

	list = r.list[r.index:]
	copy(list, list[:num-(r.maxLen-r.index)])
	return
}

//获取数据环长度
func (r *Ring)GetLen()int{
	return len(r.list)
}