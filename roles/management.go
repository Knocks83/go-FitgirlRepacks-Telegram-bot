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
)

// Search if an id corresponds to a banned user. Return index in bans list or -1 if not found.
func (roles *Roles) FindBan(id int64) int {
	for i := range roles.Blocked {
		if roles.Blocked[i] == id {
			return i
		}
	}

	return -1
}

// Ban an user from id
func (roles *Roles) AddBan(id int64) error {
	// View if current ban already exists
	if roles.FindBan(id) >= 0 {
		return nil
	}

	roles.fileMutex.Lock()

	// Update internal admin list
	roles.Blocked = append(roles.Blocked, id)

	// Generate json
	jsonBan, err := json.MarshalIndent(roles, "", "\t")
	if err != nil {
		// Delete admin just added because there is an error
		roles.Blocked = roles.Blocked[:len(roles.Blocked)-1]
		return err
	}

	// Write new json to file
	err = ioutil.WriteFile(roles.filename, jsonBan, 0644)
	if err != nil {
		// Delete admin just added because there is an I/O error
		roles.Blocked = roles.Blocked[:len(roles.Blocked)-1]
		return err
	}

	roles.fileMutex.Unlock()

	return nil
}

// Unban an user from uid
func (roles *Roles) RemoveBan(id int64) error {
	// View if current admin already exists. If no, do nothing.
	banIndex := roles.FindBan(id)
	if banIndex < 0 {
		return nil
	}

	roles.fileMutex.Lock()

	// Remove admin from slice
	roles.Blocked[banIndex] = roles.Blocked[len(roles.Blocked)-1]
	roles.Blocked = roles.Blocked[:len(roles.Blocked)-1]

	// Generate json
	jsonBan, err := json.MarshalIndent(roles, "", "\t")
	if err != nil {
		// Add admin just deleted because there is an error
		roles.Blocked = append(roles.Blocked, id)
		return err
	}

	// Write new json to file
	err = ioutil.WriteFile(roles.filename, jsonBan, 0644)
	if err != nil {
		// Add admin just deleted because there is an I/O error
		roles.Blocked = append(roles.Blocked, id)
		return err
	}

	roles.fileMutex.Unlock()

	return nil
}
