package dicom

import (
	"encoding/binary"
	"fmt"
	"log"
)

// https://www.dicomlibrary.com/dicom/transfer-syntax/

const (
	ImplicitVRLittleEndian         = "1.2.840.10008.1.2"
	ExplicitVRLittleEndian         = "1.2.840.10008.1.2.1"
	ExplicitVRBigEndian            = "1.2.840.10008.1.2.2"
	DeflatedExplicitVRLittleEndian = "1.2.840.10008.1.2.1.99"
)

// Standard list of transfer syntaxes.
var StandardTransferSyntaxes = []string{
	ImplicitVRLittleEndian,
	ExplicitVRLittleEndian,
	ExplicitVRBigEndian,
	DeflatedExplicitVRLittleEndian,
}

// Given an UID that represents a transfer syntax, return the canonical transfer
// syntax UID with the same encoding, from the list StandardTransferSyntaxes.
// Returns an error if the uid is not defined in DICOM standard, or it's not a
// transfer syntax.
//
// TODO(saito) Check the standard to see if we need to accept unknown UIDS as
// explicit little endian.
func CanonicalTransferSyntaxUID(uid string) (string, error) {
	// defaults are explicit VR, little endian
	switch uid {
	case ImplicitVRLittleEndian:
		fallthrough
	case ExplicitVRLittleEndian:
		fallthrough
	case ExplicitVRBigEndian:
		fallthrough
	case DeflatedExplicitVRLittleEndian:
		return uid, nil
	default:
		e, err := LookupUID(uid)
		if err != nil {
			return "", err
		}
		if e.Type != UIDTypeTransferSyntax {
			return "", fmt.Errorf("UID '%s' is not a transfer syntax (is %s)", uid, e.Type)
		}
		// The default is ExplicitVRLittleEndian
		return ExplicitVRLittleEndian, nil
	}
}

func ParseTransferSyntaxUID(uid string) (bo binary.ByteOrder, implicit bool, err error) {
	canonical, err := CanonicalTransferSyntaxUID(uid)
	if err != nil {
		return nil, false, err
	}
	switch canonical {
	case ImplicitVRLittleEndian:
		return binary.LittleEndian, true, nil
	case ExplicitVRLittleEndian:
		return binary.LittleEndian, false, nil
	case ExplicitVRBigEndian:
		return binary.BigEndian, false, nil
	default:
		log.Panic(canonical, uid)
		return binary.BigEndian, false, nil
	}
}
