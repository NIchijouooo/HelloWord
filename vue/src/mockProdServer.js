import { createProdMockServer } from 'vite-plugin-mock/es/createProdMockServer'
import testMock from '../mock/test'

export function setupProdMockServer() {
  createProdMockServer([...testMock])
}
