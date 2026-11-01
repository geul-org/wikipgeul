//ff:func feature=convert type=encoder control=iteration dimension=1
//ff:what NameIndex — an Issuer's documents.geul: BIND(DOC-POINTER(issuer, number, 0), docName) per document
package publish

import (
	"slices"

	"github.com/geul-org/gogeul/pkg/core"
)

// NameIndex writes the documents.geul of issuer (GEUL-ORG 문서 이름): one BIND
// per document number, Subject DOC-POINTER(issuer, number, 0), named by
// docName (q42, p31). A number without a kind letter is an error. Numbers are bound in ascending order,
// Edge Ids from 1. documents.geul has no DOC-ID (names not starting with 0x).
func NameIndex(issuer []core.Word, nums []uint64) ([]byte, error) {
	sorted := slices.Clone(nums)
	slices.Sort(sorted)
	enc := func(r core.Record) []core.Word { ws, _ := core.EncodeRecord(r); return ws }
	in := []core.Word{0x180E}
	for i, n := range slices.Compact(sorted) {
		ptr, err := core.EncodeDocPointer(issuer, []core.Word{core.Word(n >> 48), core.Word(n >> 32), core.Word(n >> 16), core.Word(n)}, []core.Word{0, 0})
		if err != nil {
			return nil, err
		}
		s, err := docName(n)
		if err != nil {
			return nil, err
		}
		name := core.EncodeBytesValue(core.CodecUTF8TextEven, []byte(s))
		p := slices.Concat(enc(core.EncodeLocalRef(uint64(i+1))), enc(ptr), enc(name))
		in = append(in, enc(core.Record{Kind: core.KindEdge, Codec: core.EdgeBind, Payload: p})...)
	}
	return core.PutWords(core.WrapDocument(append(in, 0x180F))), nil
}
