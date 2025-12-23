import Home from "./pages/Home";

export default function App() {
  return (
    /* Убираем div с bg-gray-100 и p-8. 
       Теперь Home сам контролирует свой фон и размер. */
    <Home />
  );
}