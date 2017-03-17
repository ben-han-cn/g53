package g53

import (
	"testing"
)

func TestNameConcat(t *testing.T) {
	name, _ := Root.Concat(
		NameFromStringUnsafe("a"),
		NameFromStringUnsafe("b"),
		NameFromStringUnsafe("c"),
		NameFromStringUnsafe("d"),
		NameFromStringUnsafe("e"),
		NameFromStringUnsafe("f"),
		NameFromStringUnsafe("g"),
		NameFromStringUnsafe("h"),
	)

	NameEqToStr(t, name, "a.b.c.d.e.f.g")
}

func TestNameSplit(t *testing.T) {
	wwwknetcn, _ := NewName("www.knet.Cn", true)
	n, _ := wwwknetcn.Split(0, 1)
	NameEqToStr(t, n, "www")

	n, _ = wwwknetcn.Split(0, 4)
	NameEqToStr(t, n, "www.knet.cn")

	n, _ = wwwknetcn.Split(1, 3)
	NameEqToStr(t, n, "knet.cn")

	n, _ = wwwknetcn.Split(1, 2)
	NameEqToStr(t, n, "knet.cn")

	n, _ = wwwknetcn.Parent(0)
	NameEqToStr(t, n, "www.knet.cn")

	n, _ = wwwknetcn.Parent(1)
	NameEqToStr(t, n, "knet.cn")

	n, _ = wwwknetcn.Parent(2)
	NameEqToStr(t, n, "cn")

	n, _ = wwwknetcn.Parent(3)
	NameEqToStr(t, n, ".")

	if _, err := wwwknetcn.Parent(4); err == nil {
		t.Errorf("www.knet.cn has no parent leve 4")
	}
}

func TestNameCompare(t *testing.T) {
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

	baidu_cn, _ := NewName("baidu.cn.", true)
	if cr := baidu_com.Compare(baidu_cn, true); cr.Relation != COMMONANCESTOR || cr.CommonLabelCount != 1 {
		t.Errorf("baidu.com don't have any relationship with baidu.cn", cr.Relation)
	}
}

func TestNameReverse(t *testing.T) {
	knetcn, _ := NewName("www.knet.Cn", true)
	knetcnReverse := knetcn.Reverse().String(false)
	if knetcnReverse != "cn.knet.www." {
		t.Errorf("www.knet.com reverse should be com.baidu.www. but get %v", knetcnReverse)
	}

	if Root.Reverse().String(false) != "." {
		t.Errorf("rootcom reverse should be .")
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

func TestNameHash(t *testing.T) {
	name1, _ := NewName("wwwnnnnnnnnnnnnn.KNET.cNNNNNNNNN", false)
	name2, _ := NewName("wwwnnnnnnnnnnnnn.KNET.cNNNNNNNNn", false)
	name3, _ := NewName("wwwnnnnnnnnnnnnn.KNET.cNNNNNNNNN.baidu.com.cn.net", false)
	if name1.Hash(false) != name2.Hash(false) {
		t.Errorf("same name with difference case should has same hash")
	}

	if name1.Hash(false) == name3.Hash(false) {
		t.Errorf("different name should has different hash")
	}
}

func TestNameIsSubdomain(t *testing.T) {
	www_knet_cn, _ := NewName("www.knet.Cn", true)
	www_knet, _ := NewName("www.knet", true)
	knet_cn, _ := NewName("knet.Cn", false)
	cn, _ := NewName("cn", true)
	knet, _ := NewName("kNeT", false)

	if www_knet_cn.IsSubDomain(knet_cn) == false ||
		knet_cn.IsSubDomain(cn) == false ||
		www_knet.IsSubDomain(knet) == false ||
		knet_cn.IsSubDomain(Root) == false ||
		cn.IsSubDomain(Root) == false ||
		knet.IsSubDomain(Root) == false ||
		www_knet_cn.IsSubDomain(Root) == false ||
		www_knet.IsSubDomain(Root) == false ||
		Root.IsSubDomain(Root) == false {
		t.Errorf("sub domain test fail")
	}

	if knet.IsSubDomain(knet_cn) ||
		knet.IsSubDomain(cn) ||
		Root.IsSubDomain(cn) ||
		www_knet.IsSubDomain(www_knet_cn) {
		t.Errorf("kent isnot sub domain of knet.cn or cn")
	}
}

func TestNameEquals(t *testing.T) {
	knetmixcase, _ := NewName("www.KNET.cN", false)
	knetdowncase, _ := NewName("www.knet.cn", true)
	if knetmixcase.Equals(knetdowncase) == false {
		t.Errorf("www.knet.cn is same with www.KNET.cN")
	}

	if knetmixcase.CaseSensitiveEquals(knetdowncase) {
		t.Errorf("www.knet.cn isnot casesenstive same with www.KNET.cN")
	}

	knetmixcase.Downcase()
	if knetmixcase.CaseSensitiveEquals(knetdowncase) == false {
		t.Errorf("www.knet.cn is casesenstive same with www.knet.cn")
	}
}
