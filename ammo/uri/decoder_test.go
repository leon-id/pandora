package uri

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/yandex/pandora/ammo"
)

var _ = Describe("Decoder", func() {
	It("uri decode ctx cancel", func() {
		ctx, cancel := context.WithCancel(context.Background())
		decoder := newDecoder(make(chan ammo.Ammo), ctx)
		cancel()
		err := decoder.Decode([]byte("/some/path"))
		Expect(err).To(Equal(context.Canceled))
	})
	var (
		ammoCh  chan ammo.Ammo
		decoder *decoder
	)
	BeforeEach(func() {
		ammoCh = make(chan ammo.Ammo, 10)
		decoder = newDecoder(ammoCh, context.Background())
	})
	DescribeTable("invalid input",
		func(line string) {
			err := decoder.Decode([]byte(line))
			Expect(err).NotTo(BeNil())
			Expect(ammoCh).NotTo(Receive())
			Expect(decoder.header).To(BeEmpty())
		},
		Entry("empty line", ""),
		Entry("line start", "test"),
		Entry("empty header", "[  ]"),
		Entry("no closing brace", "[key: val "),
		Entry("no header key", "[ : val ]"),
		Entry("no colon", "[ key  val ]"),
		Entry("extra space", "[key: val ] "),
	)

	Decode := func(line string) {
		err := decoder.Decode([]byte(line))
		Expect(err).To(BeNil())
	}
	It("uri", func() {
		header := http.Header{"a": []string{"b"}, "c": []string{"d"}}
		for k, v := range header {
			decoder.header[k] = v
		}
		line := "/some/path"
		Decode(line)
		var am ammo.Ammo
		Expect(ammoCh).To(Receive(&am))
		sh, ok := am.(ammo.HTTP)
		Expect(ok).To(BeTrue())
		req, sample := sh.Request()
		Expect(*req.URL).To(MatchFields(IgnoreExtras, Fields{
			"Path":   Equal(line),
			"Host":   BeEmpty(),
			"Scheme": BeEmpty(),
		}))
		Expect(req.Header).To(Equal(header))
		Expect(decoder.header).To(Equal(header))
		Expect(decoder.ammoNum).To(Equal(1))
		Expect(sample.Tags()).To(Equal("REQUEST"))
	})
	Context("header", func() {
		AfterEach(func() {
			Expect(decoder.ammoNum).To(BeZero())
		})
		It("overwrite", func() {
			decoder.header.Set("A", "b")
			Decode("[A: c]")
			Expect(decoder.header).To(Equal(http.Header{
				"A": []string{"c"},
			}))
		})
		It("add", func() {
			decoder.header.Set("A", "b")
			Decode("[C: d]")
			Expect(decoder.header).To(Equal(http.Header{
				"A": []string{"b"},
				"C": []string{"d"},
			}))
		})
		It("spaces", func() {
			Decode("[ C :   d   ]")
			Expect(decoder.header).To(Equal(http.Header{
				"C": []string{"d"},
			}))
		})
		It("value colons", func() {
			Decode("[C:c:d]")
			Expect(decoder.header).To(Equal(http.Header{
				"C": []string{"c:d"},
			}))
		})
		It("empty value", func() {
			Decode("[C:]")
			Expect(decoder.header).To(Equal(http.Header{
				"C": []string{""},
			}))
		})
	})
	It("Reset", func() {
		decoder.header.Set("a", "b")
		decoder.ammoNum = 10
		decoder.ResetHeader()
		Expect(decoder.header).To(BeEmpty())
		Expect(decoder.ammoNum).To(Equal(10))
	})

})