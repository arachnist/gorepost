// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bot

import (
	"math/rand"
	"strings"
	"time"

	"github.com/arachnist/gorepost/irc"
)

var set1 = []string{
	"To, o czym media nie chcą powiedzieć, to to, że",
	"Jak mawia kolega Michalkiewicz, to przez razwiedkę",
	"Jakby powiedział jego ekscelencja doktor Józef Goebells,",
	"Zgodnie z badaniami Instytutu Socjocybernetyki,",
	"Wbrew badaniom eurokratów, którzy robią gigantyczny przekręt na GLOBCIO,",
	"Choć Banda Czworga stara się o tym nie mówić, to",
	"Grodzkie ani - tfu! - geje nie powiedzą wam, że",
	"By oszczędzić ludzkości kolejnych konwulsji, musimy jasno i stanowczo powiedzieć naszym kobietom, starcom i dzieciom, że",
	"Żyjemy dlatego, że nie wykonano na nas aborcji, a",
	"Ja ja­koś nie miałem edu­kac­ji sek­sual­nej, a",
	"ZUS z pew­nością jest in­sty­tucją przestępczą",
	"Panie marszałku, wysoka izbo,",
	"Nawet za Hitlera czy Stalina góral mógł sobie robić oscypki jakie chciał, a dzisiaj",
	"Ja, w przeciwieństwie do papieża uważam, że",
}
var set2 = []string{
	"państwo w podatkach kradnie nam 90% dochodów",
	"na nasz koszt urządza się zawody niszczące wizerunek białego, silnego człowieka",
	"dzieci zdrowe zarażają się niepełnosprawnością od inwalidów",
	"lekka pedofilia nie jest szkodliwa społecznie",
	"generała Kiszczaka należy podziwiać",
	"wszystkie think-tanki są skażone socjalizmem",
	"kobieta potrafi zrobić z niczego trzy rzeczy: kapelusz, sałatkę i awanturę",
	"na dnie oceanu, w Rowie Mariańskim, żyją sobie ru-ru-rurkowce",
	"lewica ma jakąś obsesję na temat seksu",
	"paraolimpiada ma tyle sensu, co zawody w szachach dla debili",
	"ludzkość tak ogłupiała",
	"kobieta potrafi zrobić z niczego trzy rzeczy: kapelusz, sałatkę i awanturę",
	"System, w którym wszyscy mają pracę, nazywa się „niewolnictwo”",
	"wys­kro­bano nie tych, co trzeba",
	"związek dwóch TFU!! homosiów żadnych owoców wydać nie może",
}
var set3 = []string{
	"czego nie powstydziłby się ś.p. Adolf Hitler.",
	"o czym warto pamiętać w kontekście zmanipulowanych sondaży.",
	"podczas gdy za WCz Stalina można było swobodnie produkować oscypki.",
	"czego próbkę mieliśmy na placu Tienanmen.",
	"- ale nie będę przecież zastępować pańskich nauczycieli!",
	"choć do Murzynki nie mogłem się przemóc.",
	"a kobieta nasiąka poglądami mężczyzny przez nasienie.",
	"a czerwone jest wredne!",
	"a demokracja jest zawsze głupia.",
	", że dowolną głupotę można im wcisnąć.",
	"i trzeba być idiotą, żeby w takim ustroju żyć.",
	"natomiast socjalista myli się zawsze.",
	"żeby w małżeństwie było 50% mężczyzn i 50% kobiet",
	"bo już nie mordują, tylko kradną",
}

func korwin(output func(irc.Message), msg irc.Message) {
	if strings.Split(msg.Trailing, " ")[0] != ":korwin" {
		return
	}

	output(reply(msg, strings.Join([]string{
		set1[rand.Intn(len(set1))],
		set2[rand.Intn(len(set2))],
		set3[rand.Intn(len(set3))],
	}, " ")))

}

func init() {
	rand.Seed(time.Now().UnixNano())
	addCallback("PRIVMSG", "korwin", korwin)
}
