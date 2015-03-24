package main

import "testing"

func TestFindGlobalTrafficinHTML(t *testing.T) {
	html := []byte(`			 
	<div class="trafic">
<!--dtitre--><h5>Trafic normal sur l'ensemble des lignes de Métro.<br>
</h5><!--ftitre-->
						<!--dcontent--><!--fcontent-->			  </div>
   		should not capture this text
		</div>
	      <!--stuff--><h6>
		`)
	text := findTrafficDivText(html)
	exptext := "Trafic normal sur l'ensemble des lignes de Métro"
	if text != exptext {
		t.Errorf("got %q expected %q", text, exptext)
	}
}
