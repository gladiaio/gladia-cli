package types

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type Language string

const (
	LanguageAf   Language = "af"
	LanguageSq   Language = "sq"
	LanguageAm   Language = "am"
	LanguageAr   Language = "ar"
	LanguageHy   Language = "hy"
	LanguageAs   Language = "as"
	LanguageAz   Language = "az"
	LanguageBa   Language = "ba"
	LanguageEu   Language = "eu"
	LanguageBe   Language = "be"
	LanguageBn   Language = "bn"
	LanguageBs   Language = "bs"
	LanguageBr   Language = "br"
	LanguageBg   Language = "bg"
	LanguageCa   Language = "ca"
	LanguageZh   Language = "zh"
	LanguageHr   Language = "hr"
	LanguageCs   Language = "cs"
	LanguageDa   Language = "da"
	LanguageNl   Language = "nl"
	LanguageEn   Language = "en"
	LanguageAt   Language = "at"
	LanguageFo   Language = "fo"
	LanguageFi   Language = "fi"
	LanguageFr   Language = "fr"
	LanguageGl   Language = "gl"
	LanguageKa   Language = "ka"
	LanguageDe   Language = "de"
	LanguageEl   Language = "el"
	LanguageGu   Language = "gu"
	LanguageHt   Language = "ht"
	LanguageHa   Language = "ha"
	LanguageHaw  Language = "haw"
	LanguageHe   Language = "he"
	LanguageHi   Language = "hi"
	LanguageHu   Language = "hu"
	LanguageIs   Language = "is"
	LanguageId   Language = "id"
	LanguageIt   Language = "it"
	LanguageJp   Language = "jp"
	LanguageJv   Language = "jv"
	LanguageKn   Language = "kn"
	LanguageKk   Language = "kk"
	LanguageKm   Language = "km"
	LanguageKo   Language = "ko"
	LanguageLo   Language = "lo"
	LanguageLa   Language = "la"
	LanguageLv   Language = "lv"
	LanguageLn   Language = "ln"
	LanguageLt   Language = "lt"
	LanguageLb   Language = "lb"
	LanguageMk   Language = "mk"
	LanguageMg   Language = "mg"
	LanguageMs   Language = "ms"
	LanguageMl   Language = "ml"
	LanguageMt   Language = "mt"
	LanguageMi   Language = "mi"
	LanguageMr   Language = "mr"
	LanguageMn   Language = "mn"
	LanguageMymr Language = "mymr"
	LanguageNe   Language = "ne"
	LanguageNo   Language = "no"
	LanguageNn   Language = "nn"
	LanguageOc   Language = "oc"
	LanguagePs   Language = "ps"
	LanguageFa   Language = "fa"
	LanguagePl   Language = "pl"
	LanguagePt   Language = "pt"
	LanguagePa   Language = "pa"
	LanguageRo   Language = "ro"
	LanguageRu   Language = "ru"
	LanguageSa   Language = "sa"
	LanguageSr   Language = "sr"
	LanguageSn   Language = "sn"
	LanguageSd   Language = "sd"
	LanguageSi   Language = "si"
	LanguageSk   Language = "sk"
	LanguageSl   Language = "sl"
	LanguageSo   Language = "so"
	LanguageEs   Language = "es"
	LanguageSu   Language = "su"
	LanguageSw   Language = "sw"
	LanguageSv   Language = "sv"
	LanguageTl   Language = "tl"
	LanguageTg   Language = "tg"
	LanguageTa   Language = "ta"
	LanguageTt   Language = "tt"
	LanguageTe   Language = "te"
	LanguageTh   Language = "th"
	LanguageBo   Language = "bo"
	LanguageTr   Language = "tr"
	LanguageTk   Language = "tk"
	LanguageUk   Language = "uk"
	LanguageUr   Language = "ur"
	LanguageUz   Language = "uz"
	LanguageVi   Language = "vi"
	LanguageCy   Language = "cy"
	LanguageYi   Language = "yi"
	LanguageYo   Language = "yo"
)

type TargetLanguage string

// Constants for acceptable target languages.
const (
	TargetLanguageAf   TargetLanguage = "af"
	TargetLanguageSq   TargetLanguage = "sq"
	TargetLanguageAm   TargetLanguage = "am"
	TargetLanguageAr   TargetLanguage = "ar"
	TargetLanguageHy   TargetLanguage = "hy"
	TargetLanguageAs   TargetLanguage = "as"
	TargetLanguageAz   TargetLanguage = "az"
	TargetLanguageBa   TargetLanguage = "ba"
	TargetLanguageEu   TargetLanguage = "eu"
	TargetLanguageBe   TargetLanguage = "be"
	TargetLanguageBn   TargetLanguage = "bn"
	TargetLanguageBs   TargetLanguage = "bs"
	TargetLanguageBr   TargetLanguage = "br"
	TargetLanguageBg   TargetLanguage = "bg"
	TargetLanguageCa   TargetLanguage = "ca"
	TargetLanguageZh   TargetLanguage = "zh"
	TargetLanguageHr   TargetLanguage = "hr"
	TargetLanguageCs   TargetLanguage = "cs"
	TargetLanguageDa   TargetLanguage = "da"
	TargetLanguageNl   TargetLanguage = "nl"
	TargetLanguageEn   TargetLanguage = "en"
	TargetLanguageAt   TargetLanguage = "at"
	TargetLanguageFo   TargetLanguage = "fo"
	TargetLanguageFi   TargetLanguage = "fi"
	TargetLanguageFr   TargetLanguage = "fr"
	TargetLanguageGl   TargetLanguage = "gl"
	TargetLanguageKa   TargetLanguage = "ka"
	TargetLanguageDe   TargetLanguage = "de"
	TargetLanguageEl   TargetLanguage = "el"
	TargetLanguageGu   TargetLanguage = "gu"
	TargetLanguageHt   TargetLanguage = "ht"
	TargetLanguageHa   TargetLanguage = "ha"
	TargetLanguageHaw  TargetLanguage = "haw"
	TargetLanguageHe   TargetLanguage = "he"
	TargetLanguageHi   TargetLanguage = "hi"
	TargetLanguageHu   TargetLanguage = "hu"
	TargetLanguageIs   TargetLanguage = "is"
	TargetLanguageId   TargetLanguage = "id"
	TargetLanguageIt   TargetLanguage = "it"
	TargetLanguageJp   TargetLanguage = "jp"
	TargetLanguageJv   TargetLanguage = "jv"
	TargetLanguageKn   TargetLanguage = "kn"
	TargetLanguageKk   TargetLanguage = "kk"
	TargetLanguageKm   TargetLanguage = "km"
	TargetLanguageKo   TargetLanguage = "ko"
	TargetLanguageLo   TargetLanguage = "lo"
	TargetLanguageLa   TargetLanguage = "la"
	TargetLanguageLv   TargetLanguage = "lv"
	TargetLanguageLn   TargetLanguage = "ln"
	TargetLanguageLt   TargetLanguage = "lt"
	TargetLanguageLb   TargetLanguage = "lb"
	TargetLanguageMk   TargetLanguage = "mk"
	TargetLanguageMg   TargetLanguage = "mg"
	TargetLanguageMs   TargetLanguage = "ms"
	TargetLanguageMl   TargetLanguage = "ml"
	TargetLanguageMt   TargetLanguage = "mt"
	TargetLanguageMi   TargetLanguage = "mi"
	TargetLanguageMr   TargetLanguage = "mr"
	TargetLanguageMn   TargetLanguage = "mn"
	TargetLanguageMymr TargetLanguage = "mymr"
	TargetLanguageNe   TargetLanguage = "ne"
	TargetLanguageNo   TargetLanguage = "no"
	TargetLanguageNn   TargetLanguage = "nn"
	TargetLanguageOc   TargetLanguage = "oc"
	TargetLanguagePs   TargetLanguage = "ps"
	TargetLanguageFa   TargetLanguage = "fa"
	TargetLanguagePl   TargetLanguage = "pl"
	TargetLanguagePt   TargetLanguage = "pt"
	TargetLanguagePa   TargetLanguage = "pa"
	TargetLanguageRo   TargetLanguage = "ro"
	TargetLanguageRu   TargetLanguage = "ru"
	TargetLanguageSa   TargetLanguage = "sa"
	TargetLanguageSr   TargetLanguage = "sr"
	TargetLanguageSn   TargetLanguage = "sn"
	TargetLanguageSd   TargetLanguage = "sd"
	TargetLanguageSi   TargetLanguage = "si"
	TargetLanguageSk   TargetLanguage = "sk"
	TargetLanguageSl   TargetLanguage = "sl"
	TargetLanguageSo   TargetLanguage = "so"
	TargetLanguageEs   TargetLanguage = "es"
	TargetLanguageSu   TargetLanguage = "su"
	TargetLanguageSw   TargetLanguage = "sw"
	TargetLanguageSv   TargetLanguage = "sv"
	TargetLanguageTl   TargetLanguage = "tl"
	TargetLanguageTg   TargetLanguage = "tg"
	TargetLanguageTa   TargetLanguage = "ta"
	TargetLanguageTt   TargetLanguage = "tt"
	TargetLanguageTe   TargetLanguage = "te"
	TargetLanguageTh   TargetLanguage = "th"
	TargetLanguageBo   TargetLanguage = "bo"
	TargetLanguageTr   TargetLanguage = "tr"
	TargetLanguageTk   TargetLanguage = "tk"
	TargetLanguageUk   TargetLanguage = "uk"
	TargetLanguageUr   TargetLanguage = "ur"
	TargetLanguageUz   TargetLanguage = "uz"
	TargetLanguageVi   TargetLanguage = "vi"
	TargetLanguageCy   TargetLanguage = "cy"
	TargetLanguageYi   TargetLanguage = "yi"
	TargetLanguageYo   TargetLanguage = "yo"
)

func (l Language) String() string {
	return string(l)
}

func DisplayAllInputLanguagesNames() (string, error) {
	// Slice of all TargetLanguage constants
	allLanguages := []TargetLanguage{
		TargetLanguageAf,
		TargetLanguageSq,
		TargetLanguageAm,
		TargetLanguageAr,
		TargetLanguageHy,
		TargetLanguageAs,
		TargetLanguageAz,
		TargetLanguageBa,
		TargetLanguageEu,
		TargetLanguageBe,
		TargetLanguageBn,
		TargetLanguageBs,
		TargetLanguageBr,
		TargetLanguageBg,
		TargetLanguageCa,
		TargetLanguageZh,
		TargetLanguageHr,
		TargetLanguageCs,
		TargetLanguageDa,
		TargetLanguageNl,
		TargetLanguageEn,
		TargetLanguageAt,
		TargetLanguageFo,
		TargetLanguageFi,
		TargetLanguageFr,
		TargetLanguageGl,
		TargetLanguageKa,
		TargetLanguageDe,
		TargetLanguageEl,
		TargetLanguageGu,
		TargetLanguageHt,
		TargetLanguageHa,
		TargetLanguageHaw,
		TargetLanguageHe,
		TargetLanguageHi,
		TargetLanguageHu,
		TargetLanguageIs,
		TargetLanguageId,
		TargetLanguageIt,
		TargetLanguageJp,
		TargetLanguageJv,
		TargetLanguageKn,
		TargetLanguageKk,
		TargetLanguageKm,
		TargetLanguageKo,
		TargetLanguageLo,
		TargetLanguageLa,
		TargetLanguageLv,
		TargetLanguageLn,
		TargetLanguageLt,
		TargetLanguageLb,
		TargetLanguageMk,
		TargetLanguageMg,
		TargetLanguageMs,
		TargetLanguageMl,
		TargetLanguageMt,
		TargetLanguageMi,
		TargetLanguageMr,
		TargetLanguageMn,
		TargetLanguageNe,
		TargetLanguageNo,
		TargetLanguageNn,
		TargetLanguageOc,
		TargetLanguagePs,
		TargetLanguageFa,
		TargetLanguagePl,
		TargetLanguagePt,
		TargetLanguagePa,
		TargetLanguageRo,
		TargetLanguageRu,
		TargetLanguageSa,
		TargetLanguageSr,
		TargetLanguageSn,
		TargetLanguageSd,
		TargetLanguageSi,
		TargetLanguageSk,
		TargetLanguageSl,
		TargetLanguageSo,
		TargetLanguageEs,
		TargetLanguageSu,
		TargetLanguageSw,
		TargetLanguageSv,
		TargetLanguageTl,
		TargetLanguageTg,
		TargetLanguageTa,
		TargetLanguageTt,
		TargetLanguageTe,
		TargetLanguageTh,
		TargetLanguageBo,
		TargetLanguageTr,
		TargetLanguageTk,
		TargetLanguageUk,
		TargetLanguageUr,
		TargetLanguageUz,
		TargetLanguageVi,
		TargetLanguageCy,
		TargetLanguageYi,
		TargetLanguageYo,
	}

	for _, langCode := range allLanguages {
		tag, err := language.Parse(string(langCode))
		if err != nil {
			fmt.Printf("Error parsing language code '%s': %v\n", langCode, err)
			continue
		}
		fmt.Printf("%s: %s\n", langCode, display.English.Tags().Name(tag))
	}
	return "", nil
}

func DisplayAllTargetLanguagesNames() (string, error) {
	// Slice of all TargetLanguage constants
	allLanguages := []TargetLanguage{
		TargetLanguageAf,
		TargetLanguageSq,
		TargetLanguageAm,
		TargetLanguageAr,
		TargetLanguageHy,
		TargetLanguageAs,
		TargetLanguageAz,
		TargetLanguageBa,
		TargetLanguageEu,
		TargetLanguageBe,
		TargetLanguageBn,
		TargetLanguageBs,
		TargetLanguageBr,
		TargetLanguageBg,
		TargetLanguageCa,
		TargetLanguageZh,
		TargetLanguageHr,
		TargetLanguageCs,
		TargetLanguageDa,
		TargetLanguageNl,
		TargetLanguageEn,
		TargetLanguageAt,
		TargetLanguageFo,
		TargetLanguageFi,
		TargetLanguageFr,
		TargetLanguageGl,
		TargetLanguageKa,
		TargetLanguageDe,
		TargetLanguageEl,
		TargetLanguageGu,
		TargetLanguageHt,
		TargetLanguageHa,
		TargetLanguageHaw,
		TargetLanguageHe,
		TargetLanguageHi,
		TargetLanguageHu,
		TargetLanguageIs,
		TargetLanguageId,
		TargetLanguageIt,
		TargetLanguageJp,
		TargetLanguageJv,
		TargetLanguageKn,
		TargetLanguageKk,
		TargetLanguageKm,
		TargetLanguageKo,
		TargetLanguageLo,
		TargetLanguageLa,
		TargetLanguageLv,
		TargetLanguageLn,
		TargetLanguageLt,
		TargetLanguageLb,
		TargetLanguageMk,
		TargetLanguageMg,
		TargetLanguageMs,
		TargetLanguageMl,
		TargetLanguageMt,
		TargetLanguageMi,
		TargetLanguageMr,
		TargetLanguageMn,
		TargetLanguageNe,
		TargetLanguageNo,
		TargetLanguageNn,
		TargetLanguageOc,
		TargetLanguagePs,
		TargetLanguageFa,
		TargetLanguagePl,
		TargetLanguagePt,
		TargetLanguagePa,
		TargetLanguageRo,
		TargetLanguageRu,
		TargetLanguageSa,
		TargetLanguageSr,
		TargetLanguageSn,
		TargetLanguageSd,
		TargetLanguageSi,
		TargetLanguageSk,
		TargetLanguageSl,
		TargetLanguageSo,
		TargetLanguageEs,
		TargetLanguageSu,
		TargetLanguageSw,
		TargetLanguageSv,
		TargetLanguageTl,
		TargetLanguageTg,
		TargetLanguageTa,
		TargetLanguageTt,
		TargetLanguageTe,
		TargetLanguageTh,
		TargetLanguageBo,
		TargetLanguageTr,
		TargetLanguageTk,
		TargetLanguageUk,
		TargetLanguageUr,
		TargetLanguageUz,
		TargetLanguageVi,
		TargetLanguageCy,
		TargetLanguageYi,
		TargetLanguageYo,
	}

	for _, langCode := range allLanguages {
		tag, err := language.Parse(string(langCode))
		if err != nil {
			fmt.Printf("Error parsing language code '%s': %v\n", langCode, err)
			continue
		}
		fmt.Printf("%s: %s\n", langCode, display.English.Tags().Name(tag))
	}
	return "", nil
}
