/*
Copyright 2014 Workiva, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package palm

func mockComparator(item1, item2 interface{}) int {
	mk1 := item1.(int)
	mk2 := item2.(int)

	if mk1 == mk2 {
		return 0
	}

	if mk1 > mk2 {
		return 1
	}

	return -1
}
