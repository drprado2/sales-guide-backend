import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
    vus: 1500,
    duration: '90s',
};

export default function () {
    const url = 'http://localhost:5050/api/v1/companies/00b29a4f-2bf3-4bac-b45a-ea76605e25f8';

    const params = {
        headers: {
            'x-tenant': '01187cd0-3f54-4761-9d60-e26a005007a7',
            'x-user-id': 'bce4aab7-05c0-46a4-a0e9-baa0b4f1f35a',
            'x-email': 'Newton55@gmail.com'
        },
    };

    http.get(url, params);
    sleep(1);
}
