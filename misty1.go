// Package misty1 implements the MISTY1 cipher
/*

https://en.wikipedia.org/wiki/MISTY1
http://tools.ietf.org/search/rfc2994

*/
package misty1

import (
	"crypto/cipher"
	"strconv"
)

type misty1cipher struct {
	ek [32]uint16
}

type KeySizeError int

func (k KeySizeError) Error() string { return "misty1: invalid key size " + strconv.Itoa(int(k)) }

// New returns a cipher.Block implementing MISTY1.  The key argument must be 16 bytes.
func New(key []byte) (cipher.Block, error) {

	if l := len(key); l != 16 {
		return nil, KeySizeError(l)
	}

	m := &misty1cipher{}

	for i := 0; i < 8; i++ {
		m.ek[i] = uint16(key[i*2])*256 + uint16(key[i*2+1])
	}

	for i := 0; i < 8; i++ {
		m.ek[i+8] = fi(m.ek[i], m.ek[(i+1)%8])
		m.ek[i+16] = m.ek[i+8] & 0x1ff
		m.ek[i+24] = m.ek[i+8] >> 9
	}

	return m, nil
}

func (m *misty1cipher) BlockSize() int { return 8 }

func (m *misty1cipher) fo(fin uint32, k int) uint32 {
	t0 := uint16(fin >> 16)
	t1 := uint16(fin & 0xffff)
	t0 = t0 ^ m.ek[k]
	t0 = fi(t0, m.ek[(k+5)%8+8])
	t0 = t0 ^ t1
	t1 = t1 ^ m.ek[(k+2)%8]
	t1 = fi(t1, m.ek[(k+1)%8+8])
	t1 = t1 ^ t0
	t0 = t0 ^ m.ek[(k+7)%8]
	t0 = fi(t0, m.ek[(k+3)%8+8])
	t0 = t0 ^ t1
	t1 = t1 ^ m.ek[(k+4)%8]
	fout := uint32(t1)<<16 | uint32(t0)
	return fout
}

func fi(fin, fkey uint16) uint16 {
	d9 := fin >> 7
	d7 := fin & 0x7f
	d9 = s9table[d9] ^ d7
	d7 = s7table[d7] ^ d9
	d7 = d7 & 0x7f
	d7 = d7 ^ (fkey >> 9)
	d9 = d9 ^ (fkey & 0x1ff)
	d9 = s9table[d9] ^ d7
	fout := (d7 << 9) | d9
	return fout
}

var s7table = []uint16{
	0x1b, 0x32, 0x33, 0x5a, 0x3b, 0x10, 0x17, 0x54, 0x5b, 0x1a, 0x72, 0x73, 0x6b, 0x2c, 0x66, 0x49,
	0x1f, 0x24, 0x13, 0x6c, 0x37, 0x2e, 0x3f, 0x4a, 0x5d, 0x0f, 0x40, 0x56, 0x25, 0x51, 0x1c, 0x04,
	0x0b, 0x46, 0x20, 0x0d, 0x7b, 0x35, 0x44, 0x42, 0x2b, 0x1e, 0x41, 0x14, 0x4b, 0x79, 0x15, 0x6f,
	0x0e, 0x55, 0x09, 0x36, 0x74, 0x0c, 0x67, 0x53, 0x28, 0x0a, 0x7e, 0x38, 0x02, 0x07, 0x60, 0x29,
	0x19, 0x12, 0x65, 0x2f, 0x30, 0x39, 0x08, 0x68, 0x5f, 0x78, 0x2a, 0x4c, 0x64, 0x45, 0x75, 0x3d,
	0x59, 0x48, 0x03, 0x57, 0x7c, 0x4f, 0x62, 0x3c, 0x1d, 0x21, 0x5e, 0x27, 0x6a, 0x70, 0x4d, 0x3a,
	0x01, 0x6d, 0x6e, 0x63, 0x18, 0x77, 0x23, 0x05, 0x26, 0x76, 0x00, 0x31, 0x2d, 0x7a, 0x7f, 0x61,
	0x50, 0x22, 0x11, 0x06, 0x47, 0x16, 0x52, 0x4e, 0x71, 0x3e, 0x69, 0x43, 0x34, 0x5c, 0x58, 0x7d,
}

var s9table = []uint16{
	0x1c3, 0x0cb, 0x153, 0x19f, 0x1e3, 0x0e9, 0x0fb, 0x035, 0x181, 0x0b9, 0x117, 0x1eb, 0x133, 0x009, 0x02d, 0x0d3,
	0x0c7, 0x14a, 0x037, 0x07e, 0x0eb, 0x164, 0x193, 0x1d8, 0x0a3, 0x11e, 0x055, 0x02c, 0x01d, 0x1a2, 0x163, 0x118,
	0x14b, 0x152, 0x1d2, 0x00f, 0x02b, 0x030, 0x13a, 0x0e5, 0x111, 0x138, 0x18e, 0x063, 0x0e3, 0x0c8, 0x1f4, 0x01b,
	0x001, 0x09d, 0x0f8, 0x1a0, 0x16d, 0x1f3, 0x01c, 0x146, 0x07d, 0x0d1, 0x082, 0x1ea, 0x183, 0x12d, 0x0f4, 0x19e,
	0x1d3, 0x0dd, 0x1e2, 0x128, 0x1e0, 0x0ec, 0x059, 0x091, 0x011, 0x12f, 0x026, 0x0dc, 0x0b0, 0x18c, 0x10f, 0x1f7,
	0x0e7, 0x16c, 0x0b6, 0x0f9, 0x0d8, 0x151, 0x101, 0x14c, 0x103, 0x0b8, 0x154, 0x12b, 0x1ae, 0x017, 0x071, 0x00c,
	0x047, 0x058, 0x07f, 0x1a4, 0x134, 0x129, 0x084, 0x15d, 0x19d, 0x1b2, 0x1a3, 0x048, 0x07c, 0x051, 0x1ca, 0x023,
	0x13d, 0x1a7, 0x165, 0x03b, 0x042, 0x0da, 0x192, 0x0ce, 0x0c1, 0x06b, 0x09f, 0x1f1, 0x12c, 0x184, 0x0fa, 0x196,
	0x1e1, 0x169, 0x17d, 0x031, 0x180, 0x10a, 0x094, 0x1da, 0x186, 0x13e, 0x11c, 0x060, 0x175, 0x1cf, 0x067, 0x119,
	0x065, 0x068, 0x099, 0x150, 0x008, 0x007, 0x17c, 0x0b7, 0x024, 0x019, 0x0de, 0x127, 0x0db, 0x0e4, 0x1a9, 0x052,
	0x109, 0x090, 0x19c, 0x1c1, 0x028, 0x1b3, 0x135, 0x16a, 0x176, 0x0df, 0x1e5, 0x188, 0x0c5, 0x16e, 0x1de, 0x1b1,
	0x0c3, 0x1df, 0x036, 0x0ee, 0x1ee, 0x0f0, 0x093, 0x049, 0x09a, 0x1b6, 0x069, 0x081, 0x125, 0x00b, 0x05e, 0x0b4,
	0x149, 0x1c7, 0x174, 0x03e, 0x13b, 0x1b7, 0x08e, 0x1c6, 0x0ae, 0x010, 0x095, 0x1ef, 0x04e, 0x0f2, 0x1fd, 0x085,
	0x0fd, 0x0f6, 0x0a0, 0x16f, 0x083, 0x08a, 0x156, 0x09b, 0x13c, 0x107, 0x167, 0x098, 0x1d0, 0x1e9, 0x003, 0x1fe,
	0x0bd, 0x122, 0x089, 0x0d2, 0x18f, 0x012, 0x033, 0x06a, 0x142, 0x0ed, 0x170, 0x11b, 0x0e2, 0x14f, 0x158, 0x131,
	0x147, 0x05d, 0x113, 0x1cd, 0x079, 0x161, 0x1a5, 0x179, 0x09e, 0x1b4, 0x0cc, 0x022, 0x132, 0x01a, 0x0e8, 0x004,
	0x187, 0x1ed, 0x197, 0x039, 0x1bf, 0x1d7, 0x027, 0x18b, 0x0c6, 0x09c, 0x0d0, 0x14e, 0x06c, 0x034, 0x1f2, 0x06e,
	0x0ca, 0x025, 0x0ba, 0x191, 0x0fe, 0x013, 0x106, 0x02f, 0x1ad, 0x172, 0x1db, 0x0c0, 0x10b, 0x1d6, 0x0f5, 0x1ec,
	0x10d, 0x076, 0x114, 0x1ab, 0x075, 0x10c, 0x1e4, 0x159, 0x054, 0x11f, 0x04b, 0x0c4, 0x1be, 0x0f7, 0x029, 0x0a4,
	0x00e, 0x1f0, 0x077, 0x04d, 0x17a, 0x086, 0x08b, 0x0b3, 0x171, 0x0bf, 0x10e, 0x104, 0x097, 0x15b, 0x160, 0x168,
	0x0d7, 0x0bb, 0x066, 0x1ce, 0x0fc, 0x092, 0x1c5, 0x06f, 0x016, 0x04a, 0x0a1, 0x139, 0x0af, 0x0f1, 0x190, 0x00a,
	0x1aa, 0x143, 0x17b, 0x056, 0x18d, 0x166, 0x0d4, 0x1fb, 0x14d, 0x194, 0x19a, 0x087, 0x1f8, 0x123, 0x0a7, 0x1b8,
	0x141, 0x03c, 0x1f9, 0x140, 0x02a, 0x155, 0x11a, 0x1a1, 0x198, 0x0d5, 0x126, 0x1af, 0x061, 0x12e, 0x157, 0x1dc,
	0x072, 0x18a, 0x0aa, 0x096, 0x115, 0x0ef, 0x045, 0x07b, 0x08d, 0x145, 0x053, 0x05f, 0x178, 0x0b2, 0x02e, 0x020,
	0x1d5, 0x03f, 0x1c9, 0x1e7, 0x1ac, 0x044, 0x038, 0x014, 0x0b1, 0x16b, 0x0ab, 0x0b5, 0x05a, 0x182, 0x1c8, 0x1d4,
	0x018, 0x177, 0x064, 0x0cf, 0x06d, 0x100, 0x199, 0x130, 0x15a, 0x005, 0x120, 0x1bb, 0x1bd, 0x0e0, 0x04f, 0x0d6,
	0x13f, 0x1c4, 0x12a, 0x015, 0x006, 0x0ff, 0x19b, 0x0a6, 0x043, 0x088, 0x050, 0x15f, 0x1e8, 0x121, 0x073, 0x17e,
	0x0bc, 0x0c2, 0x0c9, 0x173, 0x189, 0x1f5, 0x074, 0x1cc, 0x1e6, 0x1a8, 0x195, 0x01f, 0x041, 0x00d, 0x1ba, 0x032,
	0x03d, 0x1d1, 0x080, 0x0a8, 0x057, 0x1b9, 0x162, 0x148, 0x0d9, 0x105, 0x062, 0x07a, 0x021, 0x1ff, 0x112, 0x108,
	0x1c0, 0x0a9, 0x11d, 0x1b0, 0x1a6, 0x0cd, 0x0f3, 0x05c, 0x102, 0x05b, 0x1d9, 0x144, 0x1f6, 0x0ad, 0x0a5, 0x03a,
	0x1cb, 0x136, 0x17f, 0x046, 0x0e1, 0x01e, 0x1dd, 0x0e6, 0x137, 0x1fa, 0x185, 0x08c, 0x08f, 0x040, 0x1b5, 0x0be,
	0x078, 0x000, 0x0ac, 0x110, 0x15e, 0x124, 0x002, 0x1bc, 0x0a2, 0x0ea, 0x070, 0x1fc, 0x116, 0x15c, 0x04c, 0x1c2,
}

func (m *misty1cipher) fl(fin uint32, k int) uint32 {
	d0 := uint16(fin >> 16)
	d1 := uint16(fin & 0xffff)
	if k&1 == 0 {
		d1 = d1 ^ (d0 & m.ek[k/2])
		d0 = d0 ^ (d1 | m.ek[(k/2+6)%8+8])
	} else {
		d1 = d1 ^ (d0 & m.ek[((k-1)/2+2)%8+8])
		d0 = d0 ^ (d1 | m.ek[((k-1)/2+4)%8])
	}
	fout := uint32(d0)<<16 | uint32(d1)
	return fout
}

func (m *misty1cipher) flinv(fin uint32, k int) uint32 {
	d0 := uint16(fin >> 16)
	d1 := uint16(fin & 0xffff)
	if k&1 == 0 {
		d0 = d0 ^ (d1 | m.ek[(k/2+6)%8+8])
		d1 = d1 ^ (d0 & m.ek[k/2])
	} else {
		d0 = d0 ^ (d1 | m.ek[((k-1)/2+4)%8])
		d1 = d1 ^ (d0 & m.ek[((k-1)/2+2)%8+8])
	}
	fout := uint32(d0)<<16 | uint32(d1)
	return fout
}

func (m *misty1cipher) Encrypt(dst, src []byte) {

	d0 := getUint32(src)
	d1 := getUint32(src[4:])

	// 0 round
	d0 = m.fl(d0, 0)
	d1 = m.fl(d1, 1)
	d1 = d1 ^ m.fo(d0, 0)

	// 1 round
	d0 = d0 ^ m.fo(d1, 1)
	// 2 round
	d0 = m.fl(d0, 2)
	d1 = m.fl(d1, 3)
	d1 = d1 ^ m.fo(d0, 2)
	// 3 round
	d0 = d0 ^ m.fo(d1, 3)

	// 4 round
	d0 = m.fl(d0, 4)
	d1 = m.fl(d1, 5)
	d1 = d1 ^ m.fo(d0, 4)
	// 5 round
	d0 = d0 ^ m.fo(d1, 5)
	// 6 round
	d0 = m.fl(d0, 6)
	d1 = m.fl(d1, 7)
	d1 = d1 ^ m.fo(d0, 6)
	// 7 round
	d0 = d0 ^ m.fo(d1, 7)
	// final
	d0 = m.fl(d0, 8)
	d1 = m.fl(d1, 9)

	putUint32(dst, d1)
	putUint32(dst[4:], d0)
}

func (m *misty1cipher) Decrypt(dst, src []byte) {

	d1 := getUint32(src)
	d0 := getUint32(src[4:])

	d0 = m.flinv(d0, 8)
	d1 = m.flinv(d1, 9)
	d0 = d0 ^ m.fo(d1, 7)
	d1 = d1 ^ m.fo(d0, 6)
	d0 = m.flinv(d0, 6)
	d1 = m.flinv(d1, 7)
	d0 = d0 ^ m.fo(d1, 5)
	d1 = d1 ^ m.fo(d0, 4)
	d0 = m.flinv(d0, 4)
	d1 = m.flinv(d1, 5)

	d0 = d0 ^ m.fo(d1, 3)
	d1 = d1 ^ m.fo(d0, 2)
	d0 = m.flinv(d0, 2)
	d1 = m.flinv(d1, 3)
	d0 = d0 ^ m.fo(d1, 1)
	d1 = d1 ^ m.fo(d0, 0)
	d0 = m.flinv(d0, 0)
	d1 = m.flinv(d1, 1)

	putUint32(dst, d0)
	putUint32(dst[4:], d1)
}

// big-endian

func getUint32(b []byte) uint32 {
	return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
}

func putUint32(b []byte, v uint32) {
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
}
