import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 100 },  // Rampa até 100 usuários em 30s
    { duration: '1m', target: 100 },   // Mantém 100 usuários por 1 minuto
    { duration: '30s', target: 0 },    // Finaliza o teste (rampa down)
  ],
};

export default function () {
  const url = 'http://localhost/v1/chargebacks';
  const userId = `user_${Math.floor(Math.random() * 1000000)}`;          // Exemplo: user_1234
  const txnId = `txn_${Math.floor(Math.random() * 100000000)}`;          // Exemplo: txn_987654

  const payload = JSON.stringify({
    user_id: userId,
    transaction_id: txnId,
    reason: `Serviço não prestado na transação ${txnId}`
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const res = http.post(url, payload, params);

  check(res, {
    'status is 202 (Created)': (r) => r.status === 202,
    'status is 200 (Already exists)': (r) => r.status === 200,
  });

  // Opcional: para não sobrecarregar SEM PAUSA
  sleep(0.5);
}
