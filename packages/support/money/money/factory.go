package money

import (
	"github.com/gollin/packages/support/money/currency"
)

// defaultManager is the default singleton instance of Manager used by the factory methods.
var defaultManager = NewManager()

// FromUSD creates a new Money instance with US Dollars
func FromUSD(amount int64) *Money {
	return defaultManager.Create(amount, currency.USD)
}

// FromEUR creates a new Money instance with Euros
func FromEUR(amount int64) *Money {
	return defaultManager.Create(amount, currency.EUR)
}

// FromGBP creates a new Money instance with British Pounds
func FromGBP(amount int64) *Money {
	return defaultManager.Create(amount, currency.GBP)
}

// FromJPY creates a new Money instance with Japanese Yen
func FromJPY(amount int64) *Money {
	return defaultManager.Create(amount, currency.JPY)
}

// FromCHF creates a new Money instance with Swiss Francs
func FromCHF(amount int64) *Money {
	return defaultManager.Create(amount, currency.CHF)
}

// FromCAD creates a new Money instance with Canadian Dollars
func FromCAD(amount int64) *Money {
	return defaultManager.Create(amount, currency.CAD)
}

// FromAUD creates a new Money instance with Australian Dollars
func FromAUD(amount int64) *Money {
	return defaultManager.Create(amount, currency.AUD)
}

// FromCNY creates a new Money instance with Chinese Yuan
func FromCNY(amount int64) *Money {
	return defaultManager.Create(amount, currency.CNY)
}

// FromINR creates a new Money instance with Indian Rupees
func FromINR(amount int64) *Money {
	return defaultManager.Create(amount, currency.INR)
}

// FromBRL creates a new Money instance with Brazilian Reais
func FromBRL(amount int64) *Money {
	return defaultManager.Create(amount, currency.BRL)
}

// FromMXN creates a new Money instance with Mexican Pesos
func FromMXN(amount int64) *Money {
	return defaultManager.Create(amount, currency.MXN)
}

// FromNZD creates a new Money instance with New Zealand Dollars
func FromNZD(amount int64) *Money {
	return defaultManager.Create(amount, currency.NZD)
}

// FromSGD creates a new Money instance with Singapore Dollars
func FromSGD(amount int64) *Money {
	return defaultManager.Create(amount, currency.SGD)
}

// FromHKD creates a new Money instance with Hong Kong Dollars
func FromHKD(amount int64) *Money {
	return defaultManager.Create(amount, currency.HKD)
}

// FromNOK creates a new Money instance with Norwegian Kroner
func FromNOK(amount int64) *Money {
	return defaultManager.Create(amount, currency.NOK)
}

// FromSEK creates a new Money instance with Swedish Kronor
func FromSEK(amount int64) *Money {
	return defaultManager.Create(amount, currency.SEK)
}

// FromDKK creates a new Money instance with Danish Kroner
func FromDKK(amount int64) *Money {
	return defaultManager.Create(amount, currency.DKK)
}

// FromRUB creates a new Money instance with Russian Rubles
func FromRUB(amount int64) *Money {
	return defaultManager.Create(amount, currency.RUB)
}

// FromZAR creates a new Money instance with South African Rand
func FromZAR(amount int64) *Money {
	return defaultManager.Create(amount, currency.ZAR)
}

// FromTRY creates a new Money instance with Turkish Lira
func FromTRY(amount int64) *Money {
	return defaultManager.Create(amount, currency.TRY)
}

// FromKRW creates a new Money instance with South Korean Won
func FromKRW(amount int64) *Money {
	return defaultManager.Create(amount, currency.KRW)
}

// FromPLN creates a new Money instance with Polish Zloty
func FromPLN(amount int64) *Money {
	return defaultManager.Create(amount, currency.PLN)
}

// FromTHB creates a new Money instance with Thai Baht
func FromTHB(amount int64) *Money {
	return defaultManager.Create(amount, currency.THB)
}

// FromIDR creates a new Money instance with Indonesian Rupiah
func FromIDR(amount int64) *Money {
	return defaultManager.Create(amount, currency.IDR)
}

// FromMYR creates a new Money instance with Malaysian Ringgit
func FromMYR(amount int64) *Money {
	return defaultManager.Create(amount, currency.MYR)
}

// FromPHP creates a new Money instance with Philippine Pesos
func FromPHP(amount int64) *Money {
	return defaultManager.Create(amount, currency.PHP)
}

// FromCZK creates a new Money instance with Czech Koruny
func FromCZK(amount int64) *Money {
	return defaultManager.Create(amount, currency.CZK)
}

// FromHUF creates a new Money instance with Hungarian Forints
func FromHUF(amount int64) *Money {
	return defaultManager.Create(amount, currency.HUF)
}

// FromILS creates a new Money instance with Israeli NewExchange Shekels
func FromILS(amount int64) *Money {
	return defaultManager.Create(amount, currency.ILS)
}

// FromCLP creates a new Money instance with Chilean Pesos
func FromCLP(amount int64) *Money {
	return defaultManager.Create(amount, currency.CLP)
}

// FromAED creates a new Money instance with the United Arab Emirates Dirhams
func FromAED(amount int64) *Money {
	return defaultManager.Create(amount, currency.AED)
}

// FromSAR creates a new Money instance with Saudi Riyals
func FromSAR(amount int64) *Money {
	return defaultManager.Create(amount, currency.SAR)
}

// FromARS creates a new Money instance with Argentine Pesos
func FromARS(amount int64) *Money {
	return defaultManager.Create(amount, currency.ARS)
}

// FromCOP creates a new Money instance with Colombian Pesos
func FromCOP(amount int64) *Money {
	return defaultManager.Create(amount, currency.COP)
}

// FromTWD creates a new Money instance with NewMoney Taiwan Dollars
func FromTWD(amount int64) *Money {
	return defaultManager.Create(amount, currency.TWD)
}

// FromVND creates a new Money instance with Vietnamese Dong
func FromVND(amount int64) *Money {
	return defaultManager.Create(amount, currency.VND)
}

// FromAFN creates a new Money instance with AFN currency
func FromAFN(amount int64) *Money {
	return defaultManager.Create(amount, currency.AFN)
}

// FromALL creates a new Money instance with ALL currency
func FromALL(amount int64) *Money {
	return defaultManager.Create(amount, currency.ALL)
}

// FromAMD creates a new Money instance with AMD currency
func FromAMD(amount int64) *Money {
	return defaultManager.Create(amount, currency.AMD)
}

// FromANG creates a new Money instance with ANG currency
func FromANG(amount int64) *Money {
	return defaultManager.Create(amount, currency.ANG)
}

// FromAOA creates a new Money instance with AOA currency
func FromAOA(amount int64) *Money {
	return defaultManager.Create(amount, currency.AOA)
}

// FromAWG creates a new Money instance with AWG currency
func FromAWG(amount int64) *Money {
	return defaultManager.Create(amount, currency.AWG)
}

// FromAZN creates a new Money instance with AZN currency
func FromAZN(amount int64) *Money {
	return defaultManager.Create(amount, currency.AZN)
}

// FromBAM creates a new Money instance with BAM currency
func FromBAM(amount int64) *Money {
	return defaultManager.Create(amount, currency.BAM)
}

// FromBBD creates a new Money instance with BBD currency
func FromBBD(amount int64) *Money {
	return defaultManager.Create(amount, currency.BBD)
}

// FromBDT creates a new Money instance with BDT currency
func FromBDT(amount int64) *Money {
	return defaultManager.Create(amount, currency.BDT)
}

// FromBGN creates a new Money instance with BGN currency
func FromBGN(amount int64) *Money {
	return defaultManager.Create(amount, currency.BGN)
}

// FromBHD creates a new Money instance with BHD currency
func FromBHD(amount int64) *Money {
	return defaultManager.Create(amount, currency.BHD)
}

// FromBIF creates a new Money instance with BIF currency
func FromBIF(amount int64) *Money {
	return defaultManager.Create(amount, currency.BIF)
}

// FromBMD creates a new Money instance with BMD currency
func FromBMD(amount int64) *Money {
	return defaultManager.Create(amount, currency.BMD)
}

// FromBND creates a new Money instance with BND currency
func FromBND(amount int64) *Money {
	return defaultManager.Create(amount, currency.BND)
}

// FromBOB creates a new Money instance with BOB currency
func FromBOB(amount int64) *Money {
	return defaultManager.Create(amount, currency.BOB)
}

// FromBOV creates a new Money instance with BOV currency
func FromBOV(amount int64) *Money {
	return defaultManager.Create(amount, currency.BOV)
}

// FromBSD creates a new Money instance with BSD currency
func FromBSD(amount int64) *Money {
	return defaultManager.Create(amount, currency.BSD)
}

// FromBTN creates a new Money instance with BTN currency
func FromBTN(amount int64) *Money {
	return defaultManager.Create(amount, currency.BTN)
}

// FromBWP creates a new Money instance with BWP currency
func FromBWP(amount int64) *Money {
	return defaultManager.Create(amount, currency.BWP)
}

// FromBYN creates a new Money instance with BYN currency
func FromBYN(amount int64) *Money {
	return defaultManager.Create(amount, currency.BYN)
}

// FromBZD creates a new Money instance with BZD currency
func FromBZD(amount int64) *Money {
	return defaultManager.Create(amount, currency.BZD)
}

// FromCDF creates a new Money instance with CDF currency
func FromCDF(amount int64) *Money {
	return defaultManager.Create(amount, currency.CDF)
}

// FromCHE creates a new Money instance with CHE currency
func FromCHE(amount int64) *Money {
	return defaultManager.Create(amount, currency.CHE)
}

// FromCHW creates a new Money instance with CHW currency
func FromCHW(amount int64) *Money {
	return defaultManager.Create(amount, currency.CHW)
}

// FromCLF creates a new Money instance with CLF currency
func FromCLF(amount int64) *Money {
	return defaultManager.Create(amount, currency.CLF)
}

// FromCOU creates a new Money instance with COU currency
func FromCOU(amount int64) *Money {
	return defaultManager.Create(amount, currency.COU)
}

// FromCRC creates a new Money instance with CRC currency
func FromCRC(amount int64) *Money {
	return defaultManager.Create(amount, currency.CRC)
}

// FromCUC creates a new Money instance with CUC currency
func FromCUC(amount int64) *Money {
	return defaultManager.Create(amount, currency.CUC)
}

// FromCUP creates a new Money instance with CUP currency
func FromCUP(amount int64) *Money {
	return defaultManager.Create(amount, currency.CUP)
}

// FromCVE creates a new Money instance with CVE currency
func FromCVE(amount int64) *Money {
	return defaultManager.Create(amount, currency.CVE)
}

// FromDJF creates a new Money instance with DJF currency
func FromDJF(amount int64) *Money {
	return defaultManager.Create(amount, currency.DJF)
}

// FromDOP creates a new Money instance with DOP currency
func FromDOP(amount int64) *Money {
	return defaultManager.Create(amount, currency.DOP)
}

// FromDZD creates a new Money instance with DZD currency
func FromDZD(amount int64) *Money {
	return defaultManager.Create(amount, currency.DZD)
}

// FromEGP creates a new Money instance with EGP currency
func FromEGP(amount int64) *Money {
	return defaultManager.Create(amount, currency.EGP)
}

// FromERN creates a new Money instance with ERN currency
func FromERN(amount int64) *Money {
	return defaultManager.Create(amount, currency.ERN)
}

// FromETB creates a new Money instance with ETB currency
func FromETB(amount int64) *Money {
	return defaultManager.Create(amount, currency.ETB)
}

// FromFJD creates a new Money instance with FJD currency
func FromFJD(amount int64) *Money {
	return defaultManager.Create(amount, currency.FJD)
}

// FromFKP creates a new Money instance with FKP currency
func FromFKP(amount int64) *Money {
	return defaultManager.Create(amount, currency.FKP)
}

// FromGEL creates a new Money instance with GEL currency
func FromGEL(amount int64) *Money {
	return defaultManager.Create(amount, currency.GEL)
}

// FromGHS creates a new Money instance with GHS currency
func FromGHS(amount int64) *Money {
	return defaultManager.Create(amount, currency.GHS)
}

// FromGIP creates a new Money instance with GIP currency
func FromGIP(amount int64) *Money {
	return defaultManager.Create(amount, currency.GIP)
}

// FromGMD creates a new Money instance with GMD currency
func FromGMD(amount int64) *Money {
	return defaultManager.Create(amount, currency.GMD)
}

// FromGNF creates a new Money instance with GNF currency
func FromGNF(amount int64) *Money {
	return defaultManager.Create(amount, currency.GNF)
}

// FromGTQ creates a new Money instance with GTQ currency
func FromGTQ(amount int64) *Money {
	return defaultManager.Create(amount, currency.GTQ)
}

// FromGYD creates a new Money instance with GYD currency
func FromGYD(amount int64) *Money {
	return defaultManager.Create(amount, currency.GYD)
}

// FromHNL creates a new Money instance with HNL currency
func FromHNL(amount int64) *Money {
	return defaultManager.Create(amount, currency.HNL)
}

// FromHTG creates a new Money instance with HTG currency
func FromHTG(amount int64) *Money {
	return defaultManager.Create(amount, currency.HTG)
}

// FromIQD creates a new Money instance with IQD currency
func FromIQD(amount int64) *Money {
	return defaultManager.Create(amount, currency.IQD)
}

// FromIRR creates a new Money instance with IRR currency
func FromIRR(amount int64) *Money {
	return defaultManager.Create(amount, currency.IRR)
}

// FromISK creates a new Money instance with ISK currency
func FromISK(amount int64) *Money {
	return defaultManager.Create(amount, currency.ISK)
}

// FromJMD creates a new Money instance with JMD currency
func FromJMD(amount int64) *Money {
	return defaultManager.Create(amount, currency.JMD)
}

// FromJOD creates a new Money instance with JOD currency
func FromJOD(amount int64) *Money {
	return defaultManager.Create(amount, currency.JOD)
}

// FromKES creates a new Money instance with KES currency
func FromKES(amount int64) *Money {
	return defaultManager.Create(amount, currency.KES)
}

// FromKGS creates a new Money instance with KGS currency
func FromKGS(amount int64) *Money {
	return defaultManager.Create(amount, currency.KGS)
}

// FromKHR creates a new Money instance with KHR currency
func FromKHR(amount int64) *Money {
	return defaultManager.Create(amount, currency.KHR)
}

// FromKMF creates a new Money instance with KMF currency
func FromKMF(amount int64) *Money {
	return defaultManager.Create(amount, currency.KMF)
}

// FromKPW creates a new Money instance with KPW currency
func FromKPW(amount int64) *Money {
	return defaultManager.Create(amount, currency.KPW)
}

// FromKWD creates a new Money instance with KWD currency
func FromKWD(amount int64) *Money {
	return defaultManager.Create(amount, currency.KWD)
}

// FromKYD creates a new Money instance with KYD currency
func FromKYD(amount int64) *Money {
	return defaultManager.Create(amount, currency.KYD)
}

// FromKZT creates a new Money instance with KZT currency
func FromKZT(amount int64) *Money {
	return defaultManager.Create(amount, currency.KZT)
}

// FromLAK creates a new Money instance with LAK currency
func FromLAK(amount int64) *Money {
	return defaultManager.Create(amount, currency.LAK)
}

// FromLBP creates a new Money instance with LBP currency
func FromLBP(amount int64) *Money {
	return defaultManager.Create(amount, currency.LBP)
}

// FromLKR creates a new Money instance with LKR currency
func FromLKR(amount int64) *Money {
	return defaultManager.Create(amount, currency.LKR)
}

// FromLRD creates a new Money instance with LRD currency
func FromLRD(amount int64) *Money {
	return defaultManager.Create(amount, currency.LRD)
}

// FromLSL creates a new Money instance with LSL currency
func FromLSL(amount int64) *Money {
	return defaultManager.Create(amount, currency.LSL)
}

// FromLYD creates a new Money instance with LYD currency
func FromLYD(amount int64) *Money {
	return defaultManager.Create(amount, currency.LYD)
}

// FromMAD creates a new Money instance with MAD currency
func FromMAD(amount int64) *Money {
	return defaultManager.Create(amount, currency.MAD)
}

// FromMDL creates a new Money instance with MDL currency
func FromMDL(amount int64) *Money {
	return defaultManager.Create(amount, currency.MDL)
}

// FromMGA creates a new Money instance with MGA currency
func FromMGA(amount int64) *Money {
	return defaultManager.Create(amount, currency.MGA)
}

// FromMKD creates a new Money instance with MKD currency
func FromMKD(amount int64) *Money {
	return defaultManager.Create(amount, currency.MKD)
}

// FromMMK creates a new Money instance with MMK currency
func FromMMK(amount int64) *Money {
	return defaultManager.Create(amount, currency.MMK)
}

// FromMNT creates a new Money instance with MNT currency
func FromMNT(amount int64) *Money {
	return defaultManager.Create(amount, currency.MNT)
}

// FromMOP creates a new Money instance with MOP currency
func FromMOP(amount int64) *Money {
	return defaultManager.Create(amount, currency.MOP)
}

// FromMUR creates a new Money instance with MUR currency
func FromMUR(amount int64) *Money {
	return defaultManager.Create(amount, currency.MUR)
}

// FromMRU creates a new Money instance with MRU currency
func FromMRU(amount int64) *Money {
	return defaultManager.Create(amount, currency.MRU)
}

// FromMVR creates a new Money instance with MVR currency
func FromMVR(amount int64) *Money {
	return defaultManager.Create(amount, currency.MVR)
}

// FromMWK creates a new Money instance with MWK currency
func FromMWK(amount int64) *Money {
	return defaultManager.Create(amount, currency.MWK)
}

// FromMXV creates a new Money instance with MXV currency
func FromMXV(amount int64) *Money {
	return defaultManager.Create(amount, currency.MXV)
}

// FromMZN creates a new Money instance with MZN currency
func FromMZN(amount int64) *Money {
	return defaultManager.Create(amount, currency.MZN)
}

// FromNAD creates a new Money instance with NAD currency
func FromNAD(amount int64) *Money {
	return defaultManager.Create(amount, currency.NAD)
}

// FromNGN creates a new Money instance with NGN currency
func FromNGN(amount int64) *Money {
	return defaultManager.Create(amount, currency.NGN)
}

// FromNIO creates a new Money instance with NIO currency
func FromNIO(amount int64) *Money {
	return defaultManager.Create(amount, currency.NIO)
}

// FromNPR creates a new Money instance with NPR currency
func FromNPR(amount int64) *Money {
	return defaultManager.Create(amount, currency.NPR)
}

// FromOMR creates a new Money instance with OMR currency
func FromOMR(amount int64) *Money {
	return defaultManager.Create(amount, currency.OMR)
}

// FromPAB creates a new Money instance with PAB currency
func FromPAB(amount int64) *Money {
	return defaultManager.Create(amount, currency.PAB)
}

// FromPEN creates a new Money instance with PEN currency
func FromPEN(amount int64) *Money {
	return defaultManager.Create(amount, currency.PEN)
}

// FromPGK creates a new Money instance with PGK currency
func FromPGK(amount int64) *Money {
	return defaultManager.Create(amount, currency.PGK)
}

// FromPKR creates a new Money instance with PKR currency
func FromPKR(amount int64) *Money {
	return defaultManager.Create(amount, currency.PKR)
}

// FromPYG creates a new Money instance with PYG currency
func FromPYG(amount int64) *Money {
	return defaultManager.Create(amount, currency.PYG)
}

// FromQAR creates a new Money instance with QAR currency
func FromQAR(amount int64) *Money {
	return defaultManager.Create(amount, currency.QAR)
}

// FromRON creates a new Money instance with RON currency
func FromRON(amount int64) *Money {
	return defaultManager.Create(amount, currency.RON)
}

// FromRSD creates a new Money instance with RSD currency
func FromRSD(amount int64) *Money {
	return defaultManager.Create(amount, currency.RSD)
}

// FromRWF creates a new Money instance with RWF currency
func FromRWF(amount int64) *Money {
	return defaultManager.Create(amount, currency.RWF)
}

// FromSBD creates a new Money instance with SBD currency
func FromSBD(amount int64) *Money {
	return defaultManager.Create(amount, currency.SBD)
}

// FromSCR creates a new Money instance with SCR currency
func FromSCR(amount int64) *Money {
	return defaultManager.Create(amount, currency.SCR)
}

// FromSDG creates a new Money instance with SDG currency
func FromSDG(amount int64) *Money {
	return defaultManager.Create(amount, currency.SDG)
}

// FromSHP creates a new Money instance with SHP currency
func FromSHP(amount int64) *Money {
	return defaultManager.Create(amount, currency.SHP)
}

// FromSLE creates a new Money instance with SLE currency
func FromSLE(amount int64) *Money {
	return defaultManager.Create(amount, currency.SLE)
}

// FromSLL creates a new Money instance with SLL currency
func FromSLL(amount int64) *Money {
	return defaultManager.Create(amount, currency.SLL)
}

// FromSOS creates a new Money instance with SOS currency
func FromSOS(amount int64) *Money {
	return defaultManager.Create(amount, currency.SOS)
}

// FromSRD creates a new Money instance with SRD currency
func FromSRD(amount int64) *Money {
	return defaultManager.Create(amount, currency.SRD)
}

// FromSSP creates a new Money instance with SSP currency
func FromSSP(amount int64) *Money {
	return defaultManager.Create(amount, currency.SSP)
}

// FromSTN creates a new Money instance with STN currency
func FromSTN(amount int64) *Money {
	return defaultManager.Create(amount, currency.STN)
}

// FromSVC creates a new Money instance with SVC currency
func FromSVC(amount int64) *Money {
	return defaultManager.Create(amount, currency.SVC)
}

// FromSYP creates a new Money instance with SYP currency
func FromSYP(amount int64) *Money {
	return defaultManager.Create(amount, currency.SYP)
}

// FromSZL creates a new Money instance with SZL currency
func FromSZL(amount int64) *Money {
	return defaultManager.Create(amount, currency.SZL)
}

// FromTJS creates a new Money instance with TJS currency
func FromTJS(amount int64) *Money {
	return defaultManager.Create(amount, currency.TJS)
}

// FromTMT creates a new Money instance with TMT currency
func FromTMT(amount int64) *Money {
	return defaultManager.Create(amount, currency.TMT)
}

// FromTND creates a new Money instance with TND currency
func FromTND(amount int64) *Money {
	return defaultManager.Create(amount, currency.TND)
}

// FromTOP creates a new Money instance with TOP currency
func FromTOP(amount int64) *Money {
	return defaultManager.Create(amount, currency.TOP)
}

// FromTTD creates a new Money instance with TTD currency
func FromTTD(amount int64) *Money {
	return defaultManager.Create(amount, currency.TTD)
}

// FromTZS creates a new Money instance with TZS currency
func FromTZS(amount int64) *Money {
	return defaultManager.Create(amount, currency.TZS)
}

// FromUAH creates a new Money instance with UAH currency
func FromUAH(amount int64) *Money {
	return defaultManager.Create(amount, currency.UAH)
}

// FromUGX creates a new Money instance with UGX currency
func FromUGX(amount int64) *Money {
	return defaultManager.Create(amount, currency.UGX)
}

// FromUSN creates a new Money instance with USN currency
func FromUSN(amount int64) *Money {
	return defaultManager.Create(amount, currency.USN)
}

// FromUYI creates a new Money instance with UYI currency
func FromUYI(amount int64) *Money {
	return defaultManager.Create(amount, currency.UYI)
}

// FromUYU creates a new Money instance with UYU currency
func FromUYU(amount int64) *Money {
	return defaultManager.Create(amount, currency.UYU)
}

// FromUYW creates a new Money instance with UYW currency
func FromUYW(amount int64) *Money {
	return defaultManager.Create(amount, currency.UYW)
}

// FromUZS creates a new Money instance with UZS currency
func FromUZS(amount int64) *Money {
	return defaultManager.Create(amount, currency.UZS)
}

// FromVES creates a new Money instance with VES currency
func FromVES(amount int64) *Money {
	return defaultManager.Create(amount, currency.VES)
}

// FromVUV creates a new Money instance with VUV currency
func FromVUV(amount int64) *Money {
	return defaultManager.Create(amount, currency.VUV)
}

// FromWST creates a new Money instance with WST currency
func FromWST(amount int64) *Money {
	return defaultManager.Create(amount, currency.WST)
}

// FromXAF creates a new Money instance with XAF currency
func FromXAF(amount int64) *Money {
	return defaultManager.Create(amount, currency.XAF)
}

// FromXAG creates a new Money instance with XAG currency
func FromXAG(amount int64) *Money {
	return defaultManager.Create(amount, currency.XAG)
}

// FromXAU creates a new Money instance with XAU currency
func FromXAU(amount int64) *Money {
	return defaultManager.Create(amount, currency.XAU)
}

// FromXBA creates a new Money instance with XBA currency
func FromXBA(amount int64) *Money {
	return defaultManager.Create(amount, currency.XBA)
}

// FromXBB creates a new Money instance with XBB currency
func FromXBB(amount int64) *Money {
	return defaultManager.Create(amount, currency.XBB)
}

// FromXBC creates a new Money instance with XBC currency
func FromXBC(amount int64) *Money {
	return defaultManager.Create(amount, currency.XBC)
}

// FromXBD creates a new Money instance with XBD currency
func FromXBD(amount int64) *Money {
	return defaultManager.Create(amount, currency.XBD)
}

// FromXCD creates a new Money instance with XCD currency
func FromXCD(amount int64) *Money {
	return defaultManager.Create(amount, currency.XCD)
}

// FromXDR creates a new Money instance with XDR currency
func FromXDR(amount int64) *Money {
	return defaultManager.Create(amount, currency.XDR)
}

// FromXOF creates a new Money instance with XOF currency
func FromXOF(amount int64) *Money {
	return defaultManager.Create(amount, currency.XOF)
}

// FromXPD creates a new Money instance with XPD currency
func FromXPD(amount int64) *Money {
	return defaultManager.Create(amount, currency.XPD)
}

// FromXPF creates a new Money instance with XPF currency
func FromXPF(amount int64) *Money {
	return defaultManager.Create(amount, currency.XPF)
}

// FromXPT creates a new Money instance with XPT currency
func FromXPT(amount int64) *Money {
	return defaultManager.Create(amount, currency.XPT)
}

// FromXSU creates a new Money instance with XSU currency
func FromXSU(amount int64) *Money {
	return defaultManager.Create(amount, currency.XSU)
}

// FromXTS creates a new Money instance with XTS currency
func FromXTS(amount int64) *Money {
	return defaultManager.Create(amount, currency.XTS)
}

// FromXUA creates a new Money instance with XUA currency
func FromXUA(amount int64) *Money {
	return defaultManager.Create(amount, currency.XUA)
}

// FromXXX creates a new Money instance with XXX currency
func FromXXX(amount int64) *Money {
	return defaultManager.Create(amount, currency.XXX)
}

// FromYER creates a new Money instance with YER currency
func FromYER(amount int64) *Money {
	return defaultManager.Create(amount, currency.YER)
}

// FromZMW creates a new Money instance with ZMW currency
func FromZMW(amount int64) *Money {
	return defaultManager.Create(amount, currency.ZMW)
}

// FromZWL creates a new Money instance with ZWL currency
func FromZWL(amount int64) *Money {
	return defaultManager.Create(amount, currency.ZWL)
}
