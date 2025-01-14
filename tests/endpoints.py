import unittest
import requests as r


class MyTestCase(unittest.TestCase):
    def test_pasty_list(self):
        response = r.get('http://localhost:8000/get_pasta_list')
        self.assertEqual(response.status_code, 200)
        result = response.json()
        self.assertEqual(type(result) is list, True)
        print(f"Pasta name: {result[0]["Name"]}")


if __name__ == '__main__':
    unittest.main()
