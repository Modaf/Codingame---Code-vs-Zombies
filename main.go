package main

import "fmt"
import "math"

//fmt.Fprintln(os.Stderr, "Debug messages...")

const vitesseZombie int = 400
const vitesseAsh int = 1000
const d2AshZombie int = 2000*2000

type perso struct {
    x, y int
}

type rekt struct {
    x, y, score int
}

func d2(x, y, xx, yy int) int {
    return (x-xx)*(x-xx) + (y-yy)*(y-yy)
}

func remove(s *[]perso, i int) []perso {
    (*s)[len(*s)-1], (*s)[i] = (*s)[i], (*s)[len(*s)-1]
    return (*s)[:len(*s)-1]
}

func move(z *perso, vitesse, x, y int) {
    d := d2(z.x, z.y, x, y)
    if d < vitesse*vitesse {
        z.x = x
        z.y = y
    } else {
        t := float64(vitesse)/math.Sqrt(float64(d))
        xx := float64(z.x) + t*float64((x - z.x))
        yy := float64(z.y) + t*float64((y - z.y))
        z.x = int(math.Ceil(xx)-1)
        z.y = int(math.Ceil(yy)-1)
    }
}

func uptadeZombie(z *perso, listeHumain *[]perso, me perso) int{
    score := 0
    m := 9999999999
    var x, y int
    for _, h := range *listeHumain {
        d := d2(z.x, z.y, h.x, h.y)
        if d < m {
            m = d
            x = h.x
            y = h.y
        }
    }
    if d2(z.x, z.y, me.x, me.y) < m {
        move(z, vitesseZombie, me.x, me.y)
    } else {
        move(z, vitesseZombie, x, y)
        score -= 20
    }
    return score
}

func uptadeAllZombie(listeZombie, listeHumain *[]perso, me perso) int{
    score := 0
    for i := 0; i < len(*listeZombie); i++ {
        score += uptadeZombie(&(*listeZombie)[i], listeHumain, me) //marche :)
    }
    return score
}

func uptadeTour(me *perso, x, y int, listeZombie, listeHumain *[]perso) int{
    points := uptadeAllZombie(listeZombie, listeHumain, *me) //marche :)
    move(me, vitesseAsh, x, y)
    nbHumains := len(*listeHumain)
    nbHumains = nbHumains * nbHumains
    fibo1 := 1
    fibo2 := 1
    suppr := make([]int, 0)
    for i, z := range *listeZombie {
        if d2(z.x, z.y, me.x, me.y) < d2AshZombie {
            points = points + fibo2 * 10 * nbHumains
            tmp := fibo2
            fibo2 = fibo1 + fibo2
            fibo1 = tmp
            suppr = append(suppr, i)
        }
    }
    for j, _ := range suppr {
        i := len(*listeZombie) - j - 1
        if i < len(*listeZombie)-1 && i > -1 {
            *listeZombie = (*listeZombie)[:i+copy((*listeZombie)[i:], (*listeZombie)[i+1:])]
        }
    }
    for _, z := range *listeZombie {
        for i, h := range *listeHumain {
            if z.x == h.x && z.y == h.y {
                *listeHumain = (*listeHumain)[:i+copy((*listeHumain)[i:], (*listeHumain)[i+1:])]
            }
        }
    }
    return points
}

func recherche(me perso, listeZombie, listeHumain []perso, prof, tour int) rekt{
    if len(listeHumain) == 0 {
        return rekt{me.x, me.y, -5000}
    }
    if prof == 0 {
        //fmt.Fprintln(os.Stderr, me.x, me.y)
        return rekt{me.x, me.y, 500*len(listeHumain)}
    }
    
    mouvements := make([]perso, 4)
    
    if tour < 5 {
        for i := 0; i<4; i++ {
            if i < len(listeHumain) {
                mouvements[i] = perso{listeHumain[i].x - me.x, listeHumain[i].y - me.y}
            }
        }
    } else {
        /*mouvements[0] = perso{2000, 0}
        mouvements[1] = perso{0, 2000}
        mouvements[2] = perso{-2000, 0}
        mouvements[3] = perso{0, -2000}*/
        mouvements[0] = perso{-707, -707}
        mouvements[1] = perso{-707, 707}
        mouvements[2] = perso{707, -707}
        mouvements[3] = perso{707, 707}
    }
    
    points := -500
    retourX := listeHumain[0].x
    retourY := listeHumain[0].y
    retourZ := make([]perso, len(listeZombie))
    for _, ii := range mouvements {
        i := ii.x
        j := ii.y
        //fmt.Fprintln(os.Stderr, listeZombie)
        
        
        lz := make([]perso, len(listeZombie))
        copy(lz, listeZombie)
        lh := make([]perso, len(listeHumain))
        copy(lh, listeHumain)
        
        //memoire := len(lz)
        
        meTmp := perso {me.x, me.y}
        if me.x+i < 0 {
            i = -me.x
        }
        if me.x+i > 16000 {
            i = 16000 - me.x
        }
        if me.y+j < 0 {
            j = -me.y
        }
        if me.y+j > 16000 {
            j = 9000 - me.y
        }
        
        tmp := uptadeTour(&meTmp, me.x+i, me.y+j, &lz, &lh) //marche :)
        tmp = tmp + recherche(meTmp, lz, lh, prof-1, tour).score
        
        if tmp > points {
            points = tmp
            retourX = me.x+i
            retourY = me.y+j
            copy(retourZ, lz)
        }
    }
    return rekt{retourX, retourY, points}
                
}

func main() {
    tour := 0
    memX := 0
    memY := 0
    for {
        tour += 1
        var x, y int
        fmt.Scan(&x, &y)
        
        var humanCount int
        fmt.Scan(&humanCount)
        
        listeHumain := make ([]perso, humanCount)
        
        for i := 0; i < humanCount; i++ {
            var humanId, humanX, humanY int
            fmt.Scan(&humanId, &humanX, &humanY)
            listeHumain[i] = perso{humanX, humanY}
        }
        
        var zombieCount int
        fmt.Scan(&zombieCount)
        
        listeZombie := make ([]perso, zombieCount)
        
        for i := 0; i < zombieCount; i++ {
            var zombieId, zombieX, zombieY, zombieXNext, zombieYNext int
            fmt.Scan(&zombieId, &zombieX, &zombieY, &zombieXNext, &zombieYNext)
            listeZombie[i] = perso{zombieX, zombieY}
        }
        
        me := perso{x, y}
        prof := 8
        
        if tour == 1 {
            res := recherche(me, listeZombie, listeHumain, prof, tour)
            memX = res.x
            memY = res.y
            // fmt.Fprintln(os.Stderr, "Debug messages...")
            fmt.Println(res.x, res.y, res.score) // Your destination coordinates
        } else {
            if tour < 4 {
                fmt.Println(memX, memY, "Loop : ", 7-tour)
            } else {
                res := recherche(me, listeZombie, listeHumain, prof, tour)
                fmt.Println(res.x, res.y, res.score)
            }
        }
    }
}
