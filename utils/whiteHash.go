package utils

import "net/url"

var (
	SNEAKERS = [282]string{"shoepalace.com", "finishline.com", "jdsports.com",
		"kidsfootlocker.com", "footlocker.com", "footaction.com",
		"eastbay.com", "champssports.com", "footlocker.ca", "net-a-porter.com",
		"mrporter.com", "yeezysupply.com", "footlocker.co.uk", "footlocker.cz",
		"footlocker.pl", "footlocker.nl", "footlocker.ie", "footlocker.dk",
		"footlocker.de", "footlocker.fr", "solebox.com", "jdsports.com",
		"snipesusa.com", "nike.com", "adidas.com", "footlocker.com.au",
		"hibbett.com", "supremenewyork.com", "sivasdescalzo.com", "bstn.com",
		"onygo.com", "courir.co", "caliroots.com", "starcowparis.com",
		"slamjam.com", "off---white.com", "offspring.co.uk", "lumtest.com",
		"m.sizeofficial.ie", "m.sizeofficial.se", "m.sizeofficial.de",
		"m.sizeofficial.dk", "m.sizeofficial.fr", "size.co.uk",
		"footpatrol.com", "thehipstore.co.uk", "footdistrict.com",
		"sneakersnstuff.com", "43einhalb.com", "allikestore.com",
		"uptherestore.com", "shop.travisscott.com", "undefeated.com",
		"theclosetinc.com", "socialstatuspgh.com", "kawsone.com", "cncpts.com",
		"cactusplantfleamarket.com", "doverstreetmarket.com",
		"blendsus.com", "dtlr.com", "juicestore.com", "stussy.com",
		"deadstock.ca", "jimmyjazz.com", "culturekings.com.au", "kith.com",
		"palaceskateboards.com", "shopnicekicks.com", "a-ma-maniere.com",
		"saintalfred.com", "socialstatuspgh.com", "rsvpgallery.com",
		"shopify.com", "nordstrom.com", "hottopic.com", "yeezysupply.com",
		"supremenewyork.com", "supremenewyork.com", "hibbett.com",
		"jdsports.com", "champssports.com", "eastbay.com", "footlocker.ca",
		"lacoste.com", "100thieves.com", "12amrun.com", "1290sqm.com",
		"abovethecloudsstore.com", "addictmiami.com", "alifenewyork.com",
		"alumniofny.com", "a-ma-maniere.com", "aimeleondore.com",
		"amongstfew.com", "apbstore.com", "atlasskateboarding.com",
		"atmosny.com", "bapeonline.com", "bbbranded.com", "bbcicecream.com",
		"bdgastore.com", "www.blendsus.com", "blkmkt.us",
		"bowsandarrowsberkeley.com", "wearebraindead.com",
		"burnrubbersneakers.com", "capsuletoronto.com", "cdgcdgcdg.com",
		"centretx.com", "cityblueshop.com", "cncpts.com",
		"commonwealth-ftgg.com", "concrete.nl", "corporategotem.com",
		"courtsidesneakers.com", "crusoeandsons.com", "culturekings.com.au",
		"thedarksideinitiative.com", "deadstock.ca", "dtlr.com",
		"dope-factory.com", "doverstreetmarket.com", "exclucitylife.com",
		"extrabutterny.com", "fearofgod.com", "feature.com", "ficegallery.com",
		"freshragsfl.com", "futurarchives.com", "gbny.com", "hanon-shop.com",
		"footlocker.lu", "footlocker.be", "footlocker.it", "footlocker.es",
		"footlocker.se", "footlocker.no", "footlocker.at", "footlocker.gr",
		"footlocker.hu", "footlocker.pt", "eflash.doverstreetmarket.com",
		"cncpts.com", "blendsus.com", "dtlr.com",
		"rsvpgallery.com", "consortium.co.uk", "snipes.com",
		"eflash-us.doverstreetmarket.com", "hannibalstore.it",
		"havenshop.com", "humanmade.jp", "huntinglodge.no", "jimmyjazz.com",
		"johnelliott.com", "juicestore.com", "kith.com", "kongonline.co.uk",
		"laceupnyc.com", "lapstoneandhammer.com", "ldrs1354.com",
		"machusonline.com", "manorphx.com", "marathonsports.com",
		"noirfonce.eu", "nojokicks.com", "notre-shop.com", "nrml.ca",
		"offthehook.ca", "oipolloi.com", "row.oneblockdown.it",
		"onenessboutique.com", "packershoes.com", "patta.nl",
		"pampamlondon.com", "par5milano.com", "properlbc.com",
		"privatesneakers.com", "reigningchamp.com", "rh-ude.com",
		"rsvpgallery.com", "saintalfred.com", "shoegallerymiami.com",
		"shopnicekicks.com", "shopsizeusa.com", "sneakerjunkiesusa.com",
		"sneakerpolitics.com", "sneakerworldshop.com",
		"socialstatuspgh.com", "soleclassics.com", "solefly.com",
		"soleheaven.com", "solestop.com", "stampd.com", "stashedsf.com",
		"stay-rooted.com", "stoneisland.com", "suede-store.com",
		"thechimpstore.com", "theclosetinc.com", "thepremierstore.com",
		"us.thesportsedit.com", "thesurestore.com", "trophyroomstore.com",
		"undefeated.com", "unheardofbrand.com", "store.unionlosangeles.com",
		"unknwn.com", "urbanindustry.co.uk", "hombreofficial.com",
		"westnyc.com", "wishatl.com", "xhibition.co", "stoneisland.co.uk",
		"adyen.com", "24segons.es", "beamhill.fi", "en.afew-store.com",
		"en.afew-store.com", "basket4ballers.com", "blowoutshop.de",
		"chmielna20.pl", "courir.be", "dtlr.com", "footshop.eu",
		"empire-leshop.com", "hollywood.eu", "lockerroomstore.be",
		"shop.maha-amsterdam.com", "nakedcph.com", "opiumparis.com",
		"oqium.com", "overkillshop.com", "prodirectsoccer.com", "shinzo.paris",
		"snowbeach.com", "sotostore.com", "streetmachine.com", "stress95.com",
		"titolo.ch", "en.titoloshop.com", "tres-bien.com", "ubiqlife.com",
		"undefeated.com", "vooberlin.com", "welcomesk8.com", "woodwood.com",
		"zupport.de", "live.adyen.com", "visa.acs.cmbchina.com",
		"verifiedbyvisa.acs.touchtechpayments.com",
		"idcheck.acs.touchtechpayments.com", "patta.nlpampamlondon.com",
		"cap.attempts.securecode.com", "aacsw.3ds.verifiedbyvisa.com",
		"paypal.com", "mpsnare.iesnare.com", "courir.com", "blendsus.com",
		"lacoste.com", "blkmkt.usbowsandarrowsberkeley.com",
		"hannibalstore.ithavenshop.com", "nrml.caoffthehook.ca",
		"footlocker.my", "footlocker.sg", "stockx.com", "api.nike.com",
		"deposit.us.shopifycs.com", "snipes.com", "footlocker.eu",
		"datadome.co", "api-js.datadome.co", "geo.captcha-delivery.com",
		"prod.jdgroupmesh.cloud", "jdsports.co.uk",
		"mixpanel.com", "braintreegateway.com", "recaptcha.net", "nike.com.hk",
		"perimeterx.com", "ipinfo.io"}
	SNEAKERSMAP map[string]int
)

func init() {
	SNEAKERSMAP = make(map[string]int)
	for i, s := range SNEAKERS {
		i=1+i
		SNEAKERSMAP[s] = 0
	}
}

func GetSneakerMap(URL string) (flag bool){
	u, err := url.Parse(URL)
	if err != nil {
		return false
	}
	_, ok := SNEAKERSMAP[u.Hostname()]
	if ok {
		return true
	} else {
		return false
	}
}
