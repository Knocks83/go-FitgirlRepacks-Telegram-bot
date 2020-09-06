/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package roles

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type Roles struct {
	Blocked   []int64 `json:"blocked"`
	filename  string
	fileMutex *sync.Mutex
}

// Create a new file roles instance from filename and return Roles pointer
func NewRoles(filename string) (*Roles, error) {
	// Instantiate a new roles struct
	roles := new(Roles)
	roles.Blocked = make([]int64, 0)
	roles.fileMutex = &sync.Mutex{}

	// Read file to a byte slice
	rolesFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	roles.filename = filename

	// Decode json file to roles struct
	err = json.Unmarshal(rolesFile, roles)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// Close the connection with a previously opened roles
func (roles *Roles) Close() {
	roles.Blocked = nil
	roles.fileMutex = nil
	roles.filename = ""
}
