package libv2ray

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dev7dev/uri-to-json/pkgs/outbound"
	"github.com/dev7dev/uri-to-json/pkgs/parser"
	"github.com/dev7dev/uri-to-json/pkgs/utils"
	v2serial "github.com/xtls/xray-core/infra/conf/serial"
	"strings"
)

func TestConfig(ConfigureFileContent string) error {
	_, err := v2serial.LoadJSONConfig(strings.NewReader(ConfigureFileContent))
	return err
}

func getOutboundJSONIntended(oStr string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(oStr), "", "  ")
	if err != nil {
		return oStr
	}
	return string(out.Bytes())
}

func getVMESSUriAttrs(oStr string) string {
	data, err := base64.StdEncoding.DecodeString(oStr)
	if err != nil {
		return oStr
	}
	return string(data)
}

func IsXrayURI(config string) bool {
	scheme := utils.ParseScheme(config)
	switch scheme {
	case parser.SchemeVmess:
		return true
	case parser.SchemeVless:
		return true
	case parser.SchemeTrojan:
		return true
	case parser.SchemeSS:
		return true
	case parser.SchemeHysteria2: // [၁] Hysteria2 ကို ခွင့်ပြုလိုက်ပါ
		return true
	default:
		return false
	}
}

func GetXrayOutboundFromURI(rawURI string) string {
	scheme := utils.ParseScheme(rawURI)
	if scheme == "" {
		return ""
	}

	// [၂] URI ကို အရင် clean/parse လုပ်မယ် (parser.go ထဲက ParseRawUri ကို ခေါ်သုံးတာ ပိုကောင်းပါတယ်)
	// VMess အတွက်က သီးသန့် logic ရှိနေလို့ ထားခဲ့ပေမယ့် တခြားဟာတွေအတွက် parser.ParseRawUri သုံးနိုင်ပါတယ်
	if scheme == parser.SchemeVmess {
		baseVMESSUri := getVMESSUriAttrs(strings.Replace(rawURI, parser.SchemeVmess, "", 1))
		rawURI = parser.SchemeVmess + baseVMESSUri
	} else {
		// SS, Vless, Hysteria2 တို့အတွက် Base64 နဲ့ character တွေကို ရှင်းလင်းပြီးသား link ရအောင်ယူမယ်
		rawURI = parser.ParseRawUri(rawURI)
	}

	// [၃] Outbound Object ရယူခြင်း
	ob := outbound.GetOutbound(outbound.XrayCore, rawURI)
	if ob != nil {
		ob.Parse(rawURI)
		// JSON indent လုပ်ပြီး ပြန်ပေးမယ်
		return getOutboundJSONIntended(ob.GetOutboundStr())
	} else {
		fmt.Printf("Scheme %s is not parsable uri or not supported by this library.\n", scheme)
	}
	return ""
}
