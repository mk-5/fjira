package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrepareRichText(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should prepare rich text for rendering"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			text := `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas finibus metus odio, sed malesuada nibh fermentum ut. *Morbi nec est lorem. Pellentesque magna tortor, suscipit sed leo eget, ultrices aliquet orci. Nam condimentum augue a ante luctus aliquam.* Aliquam erat volutpat. Ut tristique metus ut dui interdum, sit amet rhoncus lorem porttitor. Nullam* ultrices finibus mauris in molestie. Lorem ipsum dolor sit amet, consectetur adipiscing elit.

_Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Vivamus finibus laoreet molestie. Sed facilisis bibendum dolor._ Aliquam nec felis neque. Fusce tincidunt maximus magna at ullamcorper. In hac habitasse platea dictumst. Proin eget tempor elit. Integer eleifend porta arcu, vitae efficitur risus facilisis vulputate. Morbi a porttitor mauris. In hac habitasse platea dictumst. +Vivamus tristique at sapien id efficitur. Ut sed ante sit amet lorem varius elementum. In quam sapien, hendrerit in dolor vel, facilisis blandit lacus.+

-Nullam sagittis erat at cursus finibus.- Sed volutpat hendrerit est, in convallis tortor laoreet et. Integer ipsum justo, lacinia a volutpat ultrices, vehicula vitae tellus. In hac habitasse platea dictumst. Morbi tempor bibendum ligula, nec gravida ex maximus quis. Cras pretium sem ipsum. Proin rutrum libero eget nisi mollis, vitae tempor dui luctus. Ut eget sem eget tellus iaculis congue non et lacus. Morbi quis enim bibendum, mollis urna sit amet, ullamcorper neque.

* Phasellus quam neque, sollicitudin nec magna sit amet, tristique vestibulum diam. 
* Sed turpis arcu, volutpat a sapien nec, sollicitudin gravida ipsum.
*  In posuere libero et interdum vehicula. 
* Mauris nec neque tincidunt, ullamcorper est at, pellentesque nulla. 
* Phasellus leo neque, vestibulum eu fermentum nec, vehicula at leo.
* Fusce eu massa vel massa tincidunt elementum sed eget lorem.
*  Maecenas porttitor vestibulum felis nec luctus.

Phasellus viverra, leo ac porttitor consectetur, mi ipsum rutrum orci, nec congue risus turpis et sem. Fusce rhoncus felis eget purus tristique euismod. Vestibulum eu condimentum leo. Ut at mi ut augue tincidunt tincidunt at at ligula. Nam molestie mi a massa mattis luctus. Ut ut ornare leo. Phasellus sit amet lectus eu ex fringilla malesuada. Sed feugiat quam at nunc - euismod, at vehicula lectus interdum. Cras urna massa, vulputate non ante sit amet, porta - dapibus neque. Ut quis dignissim turpis. Etiam aliquet posuere orci, et dignissim enim dictum ac. Duis vel malesuada purus. Nam sit amet feugiat tellus. Curabitur justo dolor, pharetra id quam non, iaculis viverra orci.

* asdasdaasdasd
`
			wanted := `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas finibus metus odio, sed malesuada nibh fermentum ut. 󴇨Morbi nec est lorem. Pellentesque magna tortor, suscipit sed leo eget, ultrices aliquet orci. Nam condimentum augue a ante luctus aliquam.󴇨 Aliquam erat volutpat. Ut tristique metus ut dui interdum, sit amet rhoncus lorem porttitor. Nullam* ultrices finibus mauris in molestie. Lorem ipsum dolor sit amet, consectetur adipiscing elit.

󴇩Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Vivamus finibus laoreet molestie. Sed facilisis bibendum dolor.󴇩 Aliquam nec felis neque. Fusce tincidunt maximus magna at ullamcorper. In hac habitasse platea dictumst. Proin eget tempor elit. Integer eleifend porta arcu, vitae efficitur risus facilisis vulputate. Morbi a porttitor mauris. In hac habitasse platea dictumst. 󴇪Vivamus tristique at sapien id efficitur. Ut sed ante sit amet lorem varius elementum. In quam sapien, hendrerit in dolor vel, facilisis blandit lacus.󴇪

󴇫Nullam sagittis erat at cursus finibus.󴇫 Sed volutpat hendrerit est, in convallis tortor laoreet et. Integer ipsum justo, lacinia a volutpat ultrices, vehicula vitae tellus. In hac habitasse platea dictumst. Morbi tempor bibendum ligula, nec gravida ex maximus quis. Cras pretium sem ipsum. Proin rutrum libero eget nisi mollis, vitae tempor dui luctus. Ut eget sem eget tellus iaculis congue non et lacus. Morbi quis enim bibendum, mollis urna sit amet, ullamcorper neque.

 -  Phasellus quam neque, sollicitudin nec magna sit amet, tristique vestibulum diam. 
 -  Sed turpis arcu, volutpat a sapien nec, sollicitudin gravida ipsum.
 -   In posuere libero et interdum vehicula. 
 -  Mauris nec neque tincidunt, ullamcorper est at, pellentesque nulla. 
 -  Phasellus leo neque, vestibulum eu fermentum nec, vehicula at leo.
 -  Fusce eu massa vel massa tincidunt elementum sed eget lorem.
 -   Maecenas porttitor vestibulum felis nec luctus.

Phasellus viverra, leo ac porttitor consectetur, mi ipsum rutrum orci, nec congue risus turpis et sem. Fusce rhoncus felis eget purus tristique euismod. Vestibulum eu condimentum leo. Ut at mi ut augue tincidunt tincidunt at at ligula. Nam molestie mi a massa mattis luctus. Ut ut ornare leo. Phasellus sit amet lectus eu ex fringilla malesuada. Sed feugiat quam at nunc - euismod, at vehicula lectus interdum. Cras urna massa, vulputate non ante sit amet, porta - dapibus neque. Ut quis dignissim turpis. Etiam aliquet posuere orci, et dignissim enim dictum ac. Duis vel malesuada purus. Nam sit amet feugiat tellus. Curabitur justo dolor, pharetra id quam non, iaculis viverra orci.

 -  asdasdaasdasd
`

			// when
			text = PrepareRichText(text)

			// then
			assert.Equal(t, wanted, text)
		})
	}
}
