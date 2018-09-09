/*
 * Copyright (C) 2018 The ZeepinChain Authors
 * This file is part of The ZeepinChain library.
 *
 * The ZeepinChain is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ZeepinChain is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ZeepinChain.  If not, see <http://www.gnu.org/licenses/>.

 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatGala(t *testing.T) {
	assert.Equal(t, "1", FormatGala(1000000000))
	assert.Equal(t, "1.1", FormatGala(1100000000))
	assert.Equal(t, "1.123456789", FormatGala(1123456789))
	assert.Equal(t, "1000000000.123456789", FormatGala(1000000000123456789))
	assert.Equal(t, "1000000000.000001", FormatGala(1000000000000001000))
	assert.Equal(t, "1000000000.000000001", FormatGala(1000000000000000001))
}

func TestParseGala(t *testing.T) {
	assert.Equal(t, uint64(1000000000), ParseGala("1"))
	assert.Equal(t, uint64(1000000000000000000), ParseGala("1000000000"))
	assert.Equal(t, uint64(1000000000123456789), ParseGala("1000000000.123456789"))
	assert.Equal(t, uint64(1000000000000000100), ParseGala("1000000000.0000001"))
	assert.Equal(t, uint64(1000000000000000001), ParseGala("1000000000.000000001"))
	assert.Equal(t, uint64(1000000000000000001), ParseGala("1000000000.000000001123"))
}

func TestFormatZpt(t *testing.T) {
	assert.Equal(t, "0", FormatZpt(0))
	assert.Equal(t, "1", FormatZpt(1))
	assert.Equal(t, "100", FormatZpt(100))
	assert.Equal(t, "1000000000", FormatZpt(1000000000))
}

func TestParseZpt(t *testing.T) {
	assert.Equal(t, uint64(0), ParseZpt("0"))
	assert.Equal(t, uint64(1), ParseZpt("1"))
	assert.Equal(t, uint64(1000), ParseZpt("1000"))
	assert.Equal(t, uint64(1000000000), ParseZpt("1000000000"))
	assert.Equal(t, uint64(1000000), ParseZpt("1000000.123"))
}
