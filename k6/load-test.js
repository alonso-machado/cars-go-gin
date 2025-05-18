import http from 'k6/http'
import { check, sleep } from 'k6'
import { Rate } from 'k6/metrics'

// Custom metric for tracking errors
const errorRate = new Rate('errors')

// Test configuration
export const options = {
  vus: 100,
  duration: '5m',
  thresholds: {
    errors: ['rate<0.01'], // Error rate should be less than 1%
    http_req_duration: ['p(99)<200'], // 99% of requests should be below 200ms
  },
}

// Test data
const BASE_URL = 'http://localhost:8080/api/v1'
const TEST_CAR = {
  name: 'Test Car',
  brand: 'Test Brand',
  manufacturing_value: 25000,
  description: 'Test Description',
}

// Helper function to generate random data
function generateRandomCar() {
  const timestamp = new Date().getTime()
  return {
    name: `Test Car ${timestamp}`,
    brand: `Brand ${Math.floor(Math.random() * 5) + 1}`,
    manufacturing_value: Math.floor(Math.random() * 100000) + 10000,
    description: `Description for car ${timestamp}`,
  }
}

// Main test function
export default function () {
  let carId
  let responses = {}

  // Create a new car
  const newCar = generateRandomCar()
  responses.create = http.post(`${BASE_URL}/cars`, JSON.stringify(newCar), {
    headers: { 'Content-Type': 'application/json' },
  })
  check(responses.create, {
    'car creation is successful': (r) => r.status === 201,
  }) || errorRate.add(1)

  if (responses.create.status === 201) {
    carId = JSON.parse(responses.create.body).id
    sleep(1)

    // Get car by ID
    responses.getById = http.get(`${BASE_URL}/cars/${carId}`)
    check(responses.getById, {
      'get car by id is successful': (r) => r.status === 200,
    }) || errorRate.add(1)
    sleep(1)

    // Get car by name
    responses.getByName = http.get(`${BASE_URL}/cars/name/${newCar.name}`)
    check(responses.getByName, {
      'get car by name is successful': (r) => r.status === 200,
    }) || errorRate.add(1)
    sleep(1)

    // Update car
    const updatedCar = {
      ...newCar,
      name: `${newCar.name} Updated`,
      manufacturing_value: newCar.manufacturing_value + 1000,
    }
    responses.update = http.put(
      `${BASE_URL}/cars/${carId}`,
      JSON.stringify(updatedCar),
      {
        headers: { 'Content-Type': 'application/json' },
      }
    )
    check(responses.update, {
      'car update is successful': (r) => r.status === 200,
    }) || errorRate.add(1)
    sleep(1)
  }

  // Get cars by brand
  responses.getByBrand = http.get(`${BASE_URL}/cars/brand/Brand%201`)
  check(responses.getByBrand, {
    'get cars by brand is successful': (r) => r.status === 200,
  }) || errorRate.add(1)
  sleep(1)

  // Get cars by price range
  responses.getByPriceRange = http.get(
    `${BASE_URL}/cars/price-range?startPrice=10000&finalPrice=50000`
  )
  check(responses.getByPriceRange, {
    'get cars by price range is successful': (r) => r.status === 200,
  }) || errorRate.add(1)
  sleep(1)

  // Get all cars (with pagination)
  responses.getAll = http.get(`${BASE_URL}/cars?page=1&pageSize=10`)
  check(responses.getAll, {
    'get all cars is successful': (r) => r.status === 200,
  }) || errorRate.add(1)
  sleep(1)

  // Delete car (if we created one)
  if (carId) {
    responses.delete = http.del(`${BASE_URL}/cars/${carId}`)
    check(responses.delete, {
      'car deletion is successful': (r) => r.status === 204,
    }) || errorRate.add(1)
    sleep(1)
  }
}
