package cchelper

import "encoding/binary"

const PacHeadLen = 4

//package = header(4byte:bodyLen) + body
//考虑+加一个效验

func pack(body []byte) []byte {
	buf := make([]byte, len(body)+PacHeadLen)
	binary.BigEndian.PutUint32(buf[:PacHeadLen], uint32(len(body)))
	copy(buf[PacHeadLen:], body)
	return buf
}

// return truncated []byte
func unpack(buffer []byte, ch chan []byte) []byte {
reUnPa:
	bufLen := len(buffer)
	var nilBf []byte
	if bufLen < PacHeadLen { // 包长度不足,等待read
		return buffer
	}

	bodyLen := int(binary.BigEndian.Uint32(buffer[:PacHeadLen])) // 消息体长度
	//fmt.Println("bodyLen : ", bodyLen)
	bufLeftLen := bufLen - PacHeadLen

	switch {
	case bodyLen > bufLeftLen: // 包长度不足,等待read
		return buffer
	case bodyLen == bufLeftLen: // 等待read
		ch <- buffer[PacHeadLen:]
		return nilBf
	case bodyLen < bufLeftLen: // 粘包-分包
		ch <- buffer[PacHeadLen : PacHeadLen+bodyLen] // stop here!!!!
		//return unpack(buffer[PacHeadLen+bodyLen:], ch) // 递归效率会不会有问题呢?
		buffer = buffer[PacHeadLen+bodyLen:]
		goto reUnPa
	}

	return nilBf
}
