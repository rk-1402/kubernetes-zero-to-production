// Capstone 3 — HPA + Cluster Autoscaler / Karpenter under load
// Run: k6 run k6-loadtest.js
import http from 'k6/http';

export const options = { vus: 200, duration: '5m' };

export default () => http.get('http://backend.shop:8080/api/items');
