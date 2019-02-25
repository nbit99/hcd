// Copyright (c) 2015-2017 The Decred developers 
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package edwards

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"testing"
)

type ConversionVector struct {
	bIn *[32]byte
}

func testConversionVectors() []ConversionVector {
	r := rand.New(rand.NewSource(12345))

	numCvs := 50
	cvs := make([]ConversionVector, numCvs, numCvs)
	for i := 0; i < numCvs; i++ {
		bIn := new([32]byte)
		for j := 0; j < fieldIntSize; j++ {
			randByte := r.Intn(255)
			bIn[j] = uint8(randByte)
		}

		// Zero out the LSB as these aren't points.
		bIn[31] = bIn[31] &^ (1 << 7)
		cvs[i] = ConversionVector{bIn}
		r.Seed(int64(i) + 12345)
	}

	return cvs
}

// Tested functions:
//   EncodedBytesToBigInt
//   BigIntToFieldElement
//   FieldElementToEncodedBytes
//   BigIntToEncodedBytes
//   FieldElementToBigInt
//   EncodedBytesToFieldElement
func TestConversion(t *testing.T) {
	encodedNumToStrIdx := 0
	encodedNumToStrSet := []string{
		"20196841024736227335511321252453997055107605473446826399550527392484145048463",
		"20196841024736227335511321252453997055107605473446826399550527392484145048463",
		"10526890602580671421601776787753907956132695365302107847727458463041939203469",
		"6521839184025590514968509670122826212005743637868913183455710032853786234178",
		"48978957157431434435050653993003125320003416822462281103290527832248666832657",
		"4159895572539518861345595596168006899568383811471075767676415528286478158582",
		"22280614820495951172946500829116952499580924613328348188366727171047918265969",
		"49395542429090173053338999290059889648703061263880840677281131742908962906754",
		"9159517027037241827427101833531648609315728545621074713171647376409899395377",
		"30949895254060682142533175644562471377866425306664441151711563845358956591877",
		"11162413473643158193091854423259142716620990972042126831718923452075405102564",
		"8539948708472331386086579042808849879787158967647310577681201357227483798947",
		"1605590018727064783929945877682029132012855793355139251523977584085427681975",
		"45186580469309644933913280157111595250915805990139915084274696253685207655424",
		"4803736594150189698897342700036388710613749955454777749426524372915807000168",
		"44458463093829343203573846112145652781935409966482963464089847457068060406337",
		"43847912391518910515232493254376125875677511655948101409582869852756683249277",
		"7838340017403777751896157586326692776607242005376544682019613792366399557611",
		"53483903297539027073515675168328032221544053033922575353152509791677608228842",
		"8446720208822718544735330441575455413785681641890226571046615899626758103790",
		"83214700816170153628187353866879750477895001825218114349379185525906520716",
		"15672705999053019183342385294509394730658741761319919506483897620686726539506",
		"39925989312700298258602067511106679021919080633629833478073970899569606194542",
		"20378045921203371758408533697893397422378430466747056625859287444767476844853",
		"25451210723231579331405466620138966041750674441314969603476575535540200960212",
		"4474627459903707299415556391961882825477222779041974598574989918037047343964",
		"7397093240697542343237087158694422274522758581812278851156490860124123245408",
		"12019146294352615811784211023130676943046729778098323001557003733075528631035",
		"51960422165605964391089845698294985752487643873998092072153023912682573760441",
		"886059532692533180504549409050628572719111768379132175843764905063645881609",
		"51640155305131076852870025015277596315200278701141279796718402079701293657554",
		"14049973075040382983080927609020468526927855121780575811033420970575201823286",
		"6162618973699638116596472807859888770132012890154772790007061668938979292597",
		"7458672262568376308203102876435523721386473912872218538048930570929817391865",
		"21810356119146516412193497200169377227548731550005100866690685032525404754145",
		"28637616708304904312050385896309707465663047274451859776085777016997372066588",
		"17327933149852898814904118257468106846644459622769679407911385618253547550216",
		"54984262956995829168914500499743164709347162376972791975486127610358768165441",
		"16921178163410497669407012157415402253995549877675753046924185083722875330991",
		"3689950386035551021608562680180812707565873669663865379733533173639772831257",
		"16181350900485859830426932463330834642322523155640093886805355615591777762742",
		"42262715559264363958874246508050082246174775716593405560693200699826798157179",
		"14651886657373545381600338631331004518080502658076702329327433389618735343278",
		"82144858215264059972770769209255916986111959236943063998079259610333179230",
		"45155784804275518621160673294227510129350899054764902000856476028707365897276",
		"14147199139872985334438251734275254607128269954454063598011316959949979253589",
		"18900384519591001201319717443640865113984641082677367663236626750932971675075",
		"20519255595244475697549353203647543281903992059278171988053755962339471460861",
		"5175042515023984125136752466547248097261391591819683691316605910218141159117",
		"23859883219144731818639907160556442618096669420651409242145922210614455254846",
	}

	encodedBytesToStrIdx := 0
	encodedBytesToStrSet := []string{
		"8f2be510838ccf15d321efde1e209c5835112aa064b80bf0f93db988c501a72c",
		"8f2be510838ccf15d321efde1e209c5835112aa064b80bf0f93db988c501a72c",
		"8d65ca50205869a4bd3407f14598358b279e06b0f45a70b4d7c701b549024617",
		"420de51414065f49d7e2552b38356a4fd6c97d0ad8a36aac2fa5dd26ec3a6b0e",
		"113f704d1c8db9af684c992e363fd239c691390075a3c2c255950e6d7b1b496c",
		"f64a182ebe40951c6fc401020071c8eaa6db5c51da80a5a9b6d25570de6a3209",
		"71e6954c33c78e00301f690f4b8a930e07caf62049d752d85af832558f614231",
		"82224d4c02e29b43e269a79fe2daa36b8f1cec2db16d4211c96ffaf1e1e2346d",
		"316971a1452b66e4b677240531c75fa9f5aef9b3047c2ac6012aacf0581a4014",
		"053bf8a898760258f2819e7c342ae51f4c72d29f84c53476b0561a651a056d44",
		"e44d768bbc11fe2abfa92b2fc802d80020ed2b75638386f6011b2921bfb3ad18",
		"a31dc5fe9a51bd03e981c8ca11c2efb456f256e7f41f34105be89cc28e70e112",
		"b7fa2979d6fad8daddfa14b974c3df7991f68e0a3db0facc1314bdba53bb8c03",
		"00e4e3cc2de076558a5d6809345ba28691f47bbfaca5e347e96b8e0dbeb2e663",
		"68160b40f48bcc421165a942a738f63fabe76f3936206bdbb9df89368dd19e0a",
		"4172e64638f257efb822293392b6b228b1130e2004402f34667cd1932f994a62",
		"7dd6f7d87323db2cf8f0b7ac99b09d856dbf84c2fd26f97426f1ce5ff709f160",
		"eb8fc35fc213de56785ac799d40815a59c8a9677ba16f914a2bf8e09de575411",
		"eaf7493a341506307cc21708a45c605028459e15da1bc0602d9e72bd06d13e76",
		"ee620e1b4d158772c26e9a5fb3d47f48db6cc1b525ad562c13bf5a8499acac12",
		"8cba3468957c10373423e9107d4934ee9869abf5dad38fa216070bcd0c192f00",
		"f2a84c68ade91a7c653554b1feaa415b3b413c6e7c9ef7c9ab4a62e5ee6fa622",
		"6ef5384adb9de11334d4f3d0f1daeaecc68133f07659340414363a4f234f4558",
		"35892c1beeaf0bcc2a7ebc07e5c1f48be766a71e2b2378ca97a2fdb2b4900d2d",
		"d4701a1f0cb47e4219042af6d8cc8301a0a54787c9fbe59a9abf0d92f1df4438",
		"5c7bed8ce5e773cadd7999d57c7518f6b6e7390c7a50635a61808df6a48ce409",
		"602b95a90734e33b78ea2131a7f9f69c6e9dc0a87aaa89c5f6d2be743d9b5a10",
		"fbc201e7c2ec5b99abd3631b5e9766782bf7e56248eb809b841d399c8198921a",
		"b987337ab349d5b1341270a9db89ade2ac4967970d71009ab35ecae3868ee072",
		"09a52e2fae39ddd3ea0587befa18417184ec841d8c9dd143c64d75afee7df501",
		"d2edf8d94e2cc8330f5bd3c6084e464c4552ed00d0441888eced6dcdc84a2b72",
		"36ea8f93e03058e6166de376d3b06883ed086c5d8d569977d9a62578c100101f",
		"b53505eb70ce7102ed8f0adb74f25c267f8f6329ad953f45bed3110432eb9f0d",
		"f912122749bb08f77263b3178b08c721fbe8ab2867bd44947a35a3527a757d10",
		"e1d01e019f3f85676f86bba3b9e340ac4061a423740dae161da580a85e393830",
		"1c7b017d6fb3fc8a9692a19df97bc24f60a26e06dd740149d8dbcf5a1051503f",
		"08b6d739385f18e7dd572ae8e4920e1a7de9a5a0a1bafb618ed86c2e3a434f26",
		"416241faabcbc8b415f1c604ef5e558a623db38ab27c18ea192cd3826ffd8f79",
		"afe1600ad685008e0fd338f321f62ac9f708f83fb783b7bc2a5abc34250c6925",
		"19fe76b3adfb481751760f19f8bd1248099612ae5c03af2b2ad517ad1a702808",
		"b6fdb43837702a59431ce6c3f95008a600749d92c57e62cadb8c79ddee51c623",
		"7b5d7d17ef878cfbd0c2aaf6f7ed36f06b4a53927373e3444880c5f86cd96f5d",
		"aee66bacbea34f086e5d874ddf3bb589e066464e7c1ff155d4da481689ac6420",
		"5eed291f384dd85a51e2850ea51599c54c1349504c88396b8e1c26220a7e2e00",
		"3cc4ce313fe7598b2eda7e62facd386fb96aa72d4a0ac4e7b0f5f259bb44d563",
		"55fb95f4c5af8cc40216460900b1e5cbdf27999b89625556142f6ba5ec07471f",
		"c3ddfbd9862ed586afc1bd3ad2f485fdd08502827dbf01db90a5af1a113dc929",
		"fd09f92722f3798a7314ba7b8bcc72aed9b36d7e613fe1c5c268349db27c5d2d",
		"cd4e6522b33196728a8cba1fca74f57a1740a0ea070ff38e04f0984463f8700b",
		"3e3f97667479347342460ad574fd40d53c66eb41ded74b0e8e624ee91f37c034",
	}

	for _, vector := range testConversionVectors() {
		// Test encoding to FE --> bytes.
		feFB := EncodedBytesToFieldElement(vector.bIn)
		feTB := FieldElementToEncodedBytes(feFB)
		cmp := bytes.Equal(vector.bIn[:], feTB[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		// Test encoding to big int --> FE --> bytes.
		big := EncodedBytesToBigInt(vector.bIn)
		fe := BigIntToFieldElement(big)
		b := FieldElementToEncodedBytes(fe)
		cmp = bytes.Equal(vector.bIn[:], b[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		// Test encoding to big int --> bytes.
		b = BigIntToEncodedBytes(big)
		cmp = bytes.Equal(vector.bIn[:], b[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		// Test encoding FE --> big int --> bytes.
		feBig := FieldElementToBigInt(fe)
		b = BigIntToEncodedBytes(feBig)
		cmp = bytes.Equal(vector.bIn[:], b[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		// Asert our results.
		encodedNumStr := encodedNumToStrSet[encodedNumToStrIdx]
		cmp = encodedNumStr == big.String()
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}
		encodedNumToStrIdx++

		// Assert our results.
		encodedBytesToStr := encodedBytesToStrSet[encodedBytesToStrIdx]
		cmp = hex.EncodeToString(vector.bIn[:]) == encodedBytesToStr
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}
		encodedBytesToStrIdx++
	}
}

func testPointConversionVectors() []ConversionVector {
	r := rand.New(rand.NewSource(54321))

	numCvs := 50
	cvs := make([]ConversionVector, numCvs, numCvs)
	for i := 0; i < numCvs; i++ {
		bIn := new([32]byte)
		for j := 0; j < fieldIntSize; j++ {
			randByte := r.Intn(255)
			bIn[j] = uint8(randByte)
		}

		cvs[i] = ConversionVector{bIn}
		r.Seed(int64(i) + 54321)
	}

	return cvs
}

// Tested functions:
//   BigIntPointToEncodedBytes
//   extendedToBigAffine
//   EncodedBytesToBigIntPoint
func TestPointConversion(t *testing.T) {
	decodedPointsIdx := 0
	decodedPointsSet := []string{
		"36342386295235510298682738805067969701306540594271578388800019131093341795154,12122921476001995645148951048614280991245620197289177635264906062452356396947",
		"9867744591134514766234409775134503735789242596068738861666234899385955184391,27007105853956697146464305861539326123893364786296641738920699361975205011161",
		"8863951197340269192337376041401902342332353795528772831140932659547447199346,50190604070244738509034813119806333022537897198145851967099673924731480446410",
		"8339559836462603978596729027145828226446692064449465621521670593473565702096,17837149456960605255763291219950880141511042645120227896055409047531143166192",
		"29728756164057466625373568465365040229184642524788769390912342380714041647917,5732439974478537809704164552372620509962380634847718181010306158958698245522",
		"5247803078767242740120044464502520510114546116496825950704226666451120706898,35795758166829409697614900371676186086497347493565110137245013968556439289317",
		"10312254716733560945112340752820087022543905959376848337636391802511748824290,8990445345100278755573999342509674697246053177021827054216677883598688530608",
		"46187074130971894401027452608625709634897048832286753671143953908951013512232,31044683256411290311212662226177158013678758257007005486276334971220422210155",
		"43298306190744565974608750041314720022353418826989995191299464722462576232438,10299806950888062109594010166977096938659134029833804788802662958447110900184",
		"35177667294466130792822128181665344626291492152827341114635509442932953019866,30531775028158411079331658342927130455577264409965572072859069219355146042988",
		"12820030077079649627381856993531600576978343439982308457445080796710826231022,7757572025084975759300023165620534427770353848949097085637555522219082093193",
		"6253646343597756905184377893871320831940235673155211271443761113605327600367,41445277093105850390510771498467355081794651394064157691875599960362497042012",
		"7883813891724379909500855090045877805808146475175203491545972088683421931929,9625412936653339950207538235022617556749363084431208115356480470265838149235",
		"30472579368937623878795928292319870424238956006962946855422794723907421586575,54614509988478775698488822370961131128905799111181845080982179012590555357105",
		"15271018800918349528867796761843695077240549771140699068489160957947655924037,329852685195849268433605954815805257694947917358254381826586710001214418409",
		"44743338070477741089168018355702098581358101466058181816723771301716592216675,35086466151837287001694223052259529322710394149012611361861324849922889687785",
		"33866163499212710724609220769105485277006407915382983357866153206133156045778,23813080868237246772939037708015963565770032179485293012321274959431634191070",
		"5348533591149396191276147118693530904680789144057002981668376686342174347878,35905097699069976581306974641271357218166789408723537731985392719990711862278",
		"40315256029575331481422442140599294662819295345876862444982990842463168911950,16360806665429931881693966663841946038395457300165162458665310950478612690181",
		"33052698717932927149892248178127960767761537730607492185123959981890223827545,21070297496372888725360520588353008051332970094494783807974486544208068132161",
		"6067545198628916284979302752626014258394465463054351058202959902440244027330,21720867950501615598443412609276369354182820943208835740764942235573960133995",
		"53831317569733935553511538902125979694724978323896080664568625466675599453432,45539408436069992094954340278746828493900701575877952467788933528780603307924",
		"19525710214064978105626553978754569373727482210041238185590766077238636555711,15855546109268858062231319627174169988469689010275059374409951265691262712273",
		"11088014282687820959360741441437463809117498610356719851510129751795115755562,15494069385862047111068018000831027592461304482419340631336168046873859407559",
		"3780609302167966527808115014787075050391511873602293238860225507873910096878,40571044426480029736683329118811835506584488408510318455853674070731785488148",
	}

	encodedPointsIdx := 0
	encodedPointsSet := []string{
		"93b705486da83fd0e864654923104b3f13e5030c198599961b40e7079554cd1a",
		"d99e220f3d1f5eb5c04e42ddbee8ae2d0037e75493946a763fe02675ef7ab5bb",
		"ca154ac6139ad2cae9ec879d94372e5af7821babe4ad7a5c22cd9d3d0de0f66e",
		"f0f885df3d8798d555f068de9df55ddf7a27f58f0d336dfd61e1cfc304786f27",
		"922d0b01767d878b154c84e01094494e07f3da6ab9bba387679b291e3072ac8c",
		"e5bdeac8692566ac8a6cf96e13770ade7d6bf5667993ae552697f69b5fae234f",
		"b0b004306bc71b81a8f5293294dd215dbe671b3ccc2f65aeec55f3b66769e013",
		"6bb2ceb3d6d898efb79b74150c7a9543d409135c9035ee14ef87e8ce04aba244",
		"d8c9ab6c8793e53e075a8f138d36a78e20ac045a7138b0af41bf67eef07bc516",
		"6cbe9bd0c6856a79551a961bdfa6ea72be67dbbb808b832138ec2c424d5f8043",
		"89922d4d6e62a2e8cab805be279c92269ac652c587d7f9de0430ceb252a12611",
		"5c669fa7f42df778231772ad6e8ab8be5de2e329c8a186b280935f1b0f32a1db",
		"739e5b71592ecec6d2d5f8fcefd9303bd2645f80a7c242a93a91e7ac68ca4795",
		"b1dff78ebe5d0ff5f88e4245276e4d3fe038ef9855e440b9c04e13a99bb7bef8",
		"e9ed5513b18f1c2c905ec6bab764eb8226223dd443799cf429e16516a4b0ba80",
		"e9221bf1a64bf2fe52784c86ae51ad4e28448f42584a7a87cc3d1490703c92cd",
		"de021753b8b47f50e2418718c094e189a05bf7711297f1b1c8749a3be4b9a534",
		"064437ea0355acebb6d3f2129dd6a8fa826b5ce92b9dbd2e02aba443ac90614f",
		"0575769b2c4afb22b3ed5af6142163d1c4289a6d86a5c07caf4ede226fe32b24",
		"41bd4469274c4545b3e2161e9e101333e9aa45bbc2a31fea0e311fb4a25d95ae",
		"6b4d91a03c7fb93f64a9db518cef7a0c0f7cc164acc2d123e6ddb0005a930530",
		"94b726790c79d5dea2e759c070aa7ca521fc7663f51736d2412487b24a64ae64",
		"d1f93eaec626854144cdf80abe20facfad72a11159aebce2b4a49482cbeb0da3",
		"c73ea149472644acbdfc18b99e6e9ac215fdea556581b183a2bbf4a61e554122",
		"14c723f67789d320bfcccc0ff2bc8495b57b1d359d5b493aaf22df43bb65b259",
	}

	curve := new(TwistedEdwardsCurve)
	curve.InitParam25519()

	for _, vector := range testPointConversionVectors() {
		x, y, err := curve.EncodedBytesToBigIntPoint(vector.bIn)
		// The random point wasn't on the curve.
		if err != nil {
			continue
		}

		yB := BigIntPointToEncodedBytes(x, y)
		cmp := bytes.Equal(vector.bIn[:], yB[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		// Assert our results.
		var buffer bytes.Buffer
		buffer.WriteString(x.String())
		buffer.WriteString(",")
		buffer.WriteString(y.String())
		localStr := buffer.String()
		decodedPoint := decodedPointsSet[decodedPointsIdx]
		cmp = localStr == decodedPoint
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}
		decodedPointsIdx++

		// Assert our results.
		encodedPoint := encodedPointsSet[encodedPointsIdx]
		cmp = hex.EncodeToString(vector.bIn[:]) == encodedPoint
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}
		encodedPointsIdx++
	}
}

func testConversionMax() []ConversionVector {
	r := rand.New(rand.NewSource(4242314))

	numCvs := 1000000
	cvs := make([]ConversionVector, numCvs, numCvs)
	for i := 0; i < numCvs; i++ {
		bIn := new([32]byte)
		for j := 0; j < fieldIntSize; j++ {
			randByte := r.Intn(255)
			bIn[j] = uint8(randByte)
		}

		// Zero out the LSB as these aren't points.
		bIn[31] = bIn[31] &^ (1 << 7)
		cvs[i] = ConversionVector{bIn}
		r.Seed(int64(i) + 12345)
	}

	return cvs
}


