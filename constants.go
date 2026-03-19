package money

// Currency code constants for the most widely-used ISO 4217 currencies.
// These are convenience aliases for the string literals accepted by New()
// and other functions that take a currency code.
//
//	price := money.New(1000, money.USD)  // same as money.New(1000, "USD")
const (
	// Americas
	ARS = "ARS" // Argentine Peso
	AUD = "AUD" // Australian Dollar
	BOB = "BOB" // Bolivian Boliviano
	BRL = "BRL" // Brazilian Real
	BSD = "BSD" // Bahamian Dollar
	BZD = "BZD" // Belize Dollar
	CAD = "CAD" // Canadian Dollar
	CLP = "CLP" // Chilean Peso
	COP = "COP" // Colombian Peso
	DOP = "DOP" // Dominican Peso
	GYD = "GYD" // Guyanese Dollar
	HTG = "HTG" // Haitian Gourde
	JMD = "JMD" // Jamaican Dollar
	MXN = "MXN" // Mexican Peso
	NZD = "NZD" // New Zealand Dollar
	PAB = "PAB" // Panamanian Balboa
	PEN = "PEN" // Peruvian Sol
	SRD = "SRD" // Surinamese Dollar
	TTD = "TTD" // Trinidad and Tobago Dollar
	USD = "USD" // United States Dollar
	UYU = "UYU" // Uruguayan Peso

	// Europe
	ALL = "ALL" // Albanian Lek
	AZN = "AZN" // Azerbaijani Manat
	BAM = "BAM" // Bosnia-Herzegovina Convertible Mark
	BGN = "BGN" // Bulgarian Lev
	BYN = "BYN" // Belarusian Ruble
	CHF = "CHF" // Swiss Franc
	CZK = "CZK" // Czech Koruna
	DKK = "DKK" // Danish Krone
	EUR = "EUR" // Euro
	GBP = "GBP" // Pound Sterling
	GEL = "GEL" // Georgian Lari
	HRK = "HRK" // Croatian Kuna
	HUF = "HUF" // Hungarian Forint
	ISK = "ISK" // Icelandic Króna
	KZT = "KZT" // Kazakhstani Tenge
	MKD = "MKD" // Macedonian Denar
	NOK = "NOK" // Norwegian Krone
	PLN = "PLN" // Polish Złoty
	RON = "RON" // Romanian Leu
	RSD = "RSD" // Serbian Dinar
	RUB = "RUB" // Russian Ruble
	SEK = "SEK" // Swedish Krona
	TRY = "TRY" // Turkish Lira
	UAH = "UAH" // Ukrainian Hryvnia
	UZS = "UZS" // Uzbekistani Soum

	// Asia-Pacific
	BDT = "BDT" // Bangladeshi Taka
	CNY = "CNY" // Chinese Yuan
	HKD = "HKD" // Hong Kong Dollar
	IDR = "IDR" // Indonesian Rupiah
	INR = "INR" // Indian Rupee
	JPY = "JPY" // Japanese Yen
	KHR = "KHR" // Cambodian Riel
	KRW = "KRW" // South Korean Won
	LAK = "LAK" // Lao Kip
	LKR = "LKR" // Sri Lankan Rupee
	MMK = "MMK" // Myanmar Kyat
	MNT = "MNT" // Mongolian Tögrög
	MYR = "MYR" // Malaysian Ringgit
	NPR = "NPR" // Nepalese Rupee
	PHP = "PHP" // Philippine Peso
	PKR = "PKR" // Pakistani Rupee
	SGD = "SGD" // Singapore Dollar
	THB = "THB" // Thai Baht
	TWD = "TWD" // New Taiwan Dollar
	VND = "VND" // Vietnamese Dong

	// Middle East
	AED = "AED" // UAE Dirham
	BHD = "BHD" // Bahraini Dinar
	ILS = "ILS" // Israeli Shekel
	IQD = "IQD" // Iraqi Dinar
	IRR = "IRR" // Iranian Rial
	JOD = "JOD" // Jordanian Dinar
	KWD = "KWD" // Kuwaiti Dinar
	OMR = "OMR" // Omani Rial
	QAR = "QAR" // Qatari Riyal
	SAR = "SAR" // Saudi Riyal

	// Africa
	DZD = "DZD" // Algerian Dinar
	EGP = "EGP" // Egyptian Pound
	ETB = "ETB" // Ethiopian Birr
	GHS = "GHS" // Ghanaian Cedi
	KES = "KES" // Kenyan Shilling
	LYD = "LYD" // Libyan Dinar
	MAD = "MAD" // Moroccan Dirham
	MUR = "MUR" // Mauritian Rupee
	NGN = "NGN" // Nigerian Naira
	TND = "TND" // Tunisian Dinar
	TZS = "TZS" // Tanzanian Shilling
	UGX = "UGX" // Ugandan Shilling
	ZAR = "ZAR" // South African Rand
	ZMW = "ZMW" // Zambian Kwacha
)
