package utils

func Average(value ...int) float64{
  size := len(value)
  sum := 0

  for _, item := range value {
    sum += item
  }
  
  return float64(sum)/float64(size)
}