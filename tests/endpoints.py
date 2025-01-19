import unittest
import requests as r


class MyTestCase(unittest.TestCase):
    def test_pasty_list(self):
        response = r.get('http://localhost:8000/get_pasta_list')
        self.assertEqual(response.status_code, 200)
        result = response.json()
        self.assertEqual(type(result) is list, True)
        self.assertEqual(type(result[0]["Tags"]) is list, True)
        print(f"Pasta name: {result[0]["Name"]}")

    def test_login(self):
        json = {"Credentials": "1", "Password": "1"}
        response = r.post('http://localhost:8000/login', json=json)
        self.assertEqual(response.status_code, 200)
        result = response.json()
        print(f"{result["token"]}")

    def test_add_record(self):
        json = {"Credentials": "lil_cock", "Password": "1"}
        response = r.post('http://localhost:8000/login', json=json)
        self.assertEqual(response.status_code, 200)
        result = response.json()
        token = result["token"]

        json = {"Name": "Test story", "Text": """text"""
                }
        response = r.post('http://localhost:8000/add_pasta', json=json, headers={'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token})
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.text, "Success")




if __name__ == '__main__':
    unittest.main()
