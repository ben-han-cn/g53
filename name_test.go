package g53

import (
	"testing"
)

func TestNameConcat(t *testing.T) {
	knetcn, _ := NewName("www.knet.Cn", true)
	knet, _ := NewName("www.knet", true)
	cn, _ := NewName("cn", true)

	knetcn2, _ := knet.Concat(cn)
	relationship := knetcn.Compare(knetcn2, true)
	if relationship.Order != 0 ||
		relationship.CommonLabelCount != 4 ||
		relationship.Relation != EQUAL {
		t.Errorf("n1 should equal to www.knet.cn,but get %v", knetcn2.String(true))
	}

	knetcnReverse := knetcn.Reverse().String(false)

	if knetcnReverse != "cn.knet.www." {
		t.Errorf("www.knet.com reverse should be com.baidu.www. but get %v", knetcnReverse)
	}

	n, _ := knetcn.Split(0, 1)
	NameEqToStr(t, n, "www")

	n, _ = knetcn.Split(0, 4)
	NameEqToStr(t, n, "www.knet.cn")

	n, _ = knetcn.Split(1, 3)
	NameEqToStr(t, n, "knet.cn")

	n, _ = knetcn.Split(1, 2)
	NameEqToStr(t, n, "knet.cn")

	n, _ = knetcn.Parent(0)
	NameEqToStr(t, n, "www.knet.cn")

	n, _ = knetcn.Parent(1)
	NameEqToStr(t, n, "knet.cn")

	n, _ = knetcn.Parent(2)
	NameEqToStr(t, n, "cn")

	n, _ = knetcn.Parent(3)
	NameEqToStr(t, n, ".")

	if _, err := knetcn.Parent(4); err == nil {
		t.Errorf("www.knet.cn has no parent leve 4")
	}

	knetmixcase, _ := NewName("www.KNET.cN", false)
	knetdowncase, _ := NewName("www.knet.cn", true)
	knetmixcase.Downcase()
	if cr := knetmixcase.Compare(knetdowncase, true); cr.Order != 0 || cr.CommonLabelCount != 4 || cr.Relation != EQUAL {
		t.Errorf("down case failed:%v", knetmixcase)
	}

	baidu_com, _ := NewName("baidu.com.", true)
	www_baidu_com, _ := NewName("www.baidu.com", true)
	if cr := baidu_com.Compare(www_baidu_com, true); cr.Relation != SUPERDOMAIN {
		t.Errorf("baidu.com is www.baidu.com's superdomain but get %v", cr.Relation)
	}
}

func TestNameStrip(t *testing.T) {
	knetmixcase, _ := NewName("www.KNET.cN", false)
	knetWithoutCN, _ := knetmixcase.StripLeft(1)
	NameEqToStr(t, knetWithoutCN, "knet.cn")

	if knetmixcase.Hash(false) == knetWithoutCN.Hash(false) {
		t.Errorf("hash should be different if name isn't same")
	}

	cn, _ := knetmixcase.StripLeft(2)
	NameEqToStr(t, cn, "cn")

	root, _ := knetmixcase.StripLeft(3)
	NameEqToStr(t, root, ".")

	knettld, _ := knetmixcase.StripRight(1)
	NameEqToStr(t, knettld, "www.knet")

	wwwtld, _ := knetmixcase.StripRight(2)
	NameEqToStr(t, wwwtld, "www")

	wwwString := wwwtld.String(true)
	if wwwString != "www" {
		t.Errorf("wwwString to string should be www but %v", wwwString)
	}

	wwwString = wwwtld.String(false)
	if wwwString != "www." {
		t.Errorf("wwwString to string should be www. but %v", wwwString)
	}

	root, _ = knetmixcase.StripRight(3)
	NameEqToStr(t, root, ".")
}
