package qsort

func quicSort(values []int, left, right int) {
    temp := values[left]
    p := left
    i, j := left, right

    for i <= j {
        for j >= p && values[j] >= temp {
            j--
        }

        if j >= p {
            values[p] = values[j]
            p = j
        }

        if values[i] <= temp && i <= p {
            i++
        }

        if i <= p {
            values[p] = values[i]
            p = i
        }
    }

    values[p] = temp
    if p-left > 1 {
        quicSort(values, left, p-1)
    }

    if right-p > 1 {
        quicSort(values, p+1, right)
    }
}

func QuickSort(values []int) {
    quicSort(values, 0, len(values)-1)
}
