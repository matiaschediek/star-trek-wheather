# Stak Trek Wheather

## Descripción

Complejo sistema informático que predice el clima del sistema solar ficticio de Star Trek. En este existen tres civilizaciones. Vulcanos, Ferengis y Betasoides. Cada civilización vive en paz en su respectivo planeta.

## Premisas

* El planeta Ferengi se desplaza con una velocidad angular de 1 grados/día en sentido horario. Su distancia con respecto al sol es de 500Km.
* El planeta Betasoide se desplaza con una velocidad angular de 3 grados/día en sentido horario. Su distancia con respecto al sol es de 2000Km.
* El planeta Vulcano se desplaza con una velocidad angular de 5 grados/día en sentido anti­horario, su distancia con respecto al sol es de 1000Km.
* Todas las órbitas son circulares. 

1. Cuando los tres planetas están alineados entre sí y a su vez alineados con respecto al sol, el sistema solar experimenta un período de sequía ("Drought").

2. Cuando los tres planetas no están alineados, forman entre sí un triángulo. Es sabido que en el momento en el que el sol se encuentra dentro del triángulo, el sistema solar experimenta un período de lluvia ("Rainy"), teniendo éste, un pico de intensidad cuando el perímetro del triángulo está en su máximo (Este máximo es en primera instancia tomado de los valores calculados en los primeros 10 años. En caso de posteriormente calcularse un valor mas alto el sistema nivelara los datos para que contemplen este "nuevo" máximo). En caso de que el sol NO se dentro del triángulo el sistema solar no presenta anormalidades por lo que se considera normal ("Normal")

3. Las condiciones óptimas de presión y temperatura ("Optimum") se dan cuando los tres planetas están alineados entre sí pero no están alineados con el sol.

## Consumir el Sistema

Este sistema calcula automáticamente los primero 10 años posteriores a su fecha de inicio de ejecución y los resguarda en memoria para futuras consultas. Sin embargo pueden consultarse días posteriores y anteriores a dicha fecha y sus datos serán almacenados en memoria junto con los anteriores.

Para consultar el sistema expone una API REST:

* Consultar por el numero de días posteriores o anteriores (utilizar números negativos) a la fecha de inicio (En caso de no encontrar el día almacenado en memoria lo calcula y guarda en memoria)

```url
http://{host}:Port/wheather/day/{numberOfDays}
```

Ejemplo de Respuesta:

```json
{
    "Wheather":"Normal",
    "day":1
}
```

* Consultar por fecha en formato yyyy-mm-dd (En caso de no encontrar el día almacenado en memoria lo calcula y guarda en memoria)

```url
http://{host}:Port/wheather/date/{yyyy}-{mm}-{dd}
```

Ejemplo de Respuesta:

```url
http://localhost:8080/wheather/date/2020-02-10
```

Suponiendo que la fecha de inicio es el 2020-02-09

```json
{
    "Wheather":"Normal",
    "day":1
}
```

* Consultar por rango de días (En caso de no encontrar los días almacenados en memoria los calcula y guarda en memoria). En caso de que el to y el from estén invertidos el sistema los corrige automáticamente.

```url
http://{host}:Port/wheather/range/{numberOfDaysFrom}/{numberOfDaysTo}
```

Ejemplo de Respuesta:

```url
http://localhost:8080/wheather/range/1/-3
```

```json
{
    "Total":5,
    "Days":[
        {
            "Wheather":"Normal",
            "day":-3
        },
        {
            "Wheather":"Normal",
            "day":-2
        },{
            "Wheather":"Normal",
            "day":-1
        },{
            "Wheather":"Drought",
            "day":0
        },
        {
            "Wheather":"Normal",
            "day":1
        }
    ]
}
```

* Consultar por tipo de clima (Solo traerá los almacenados en memoria). Siendo eso tipos: "Normal", "Drought", "Optimum", "Rainy". En caso del ultimo ("Rainy") tambien se informaran la cantidad de días que se encontró la lluvia mas intensa y una lista de esos días. 

```url
http://{host}:Port/wheather/type/{type}
```

Ejemplo de Respuesta:

```url
http://localhost:8080/wheather/type/Normal
```

```json
{
    "Total":1813,
    "Periods":61,
    "Days":[
        {
            "Wheather":"Normal",
            "day":-3
        },
        {
            "Wheather":"Normal",
            "day":-2
        },{
            "Wheather":"Normal",
            "day":-1
        },
        {
            "Wheather":"Normal",
            "day":1
        }
        ...
    ]
}
```

```url
http://localhost:8080/wheather/type/Rainy
```


```json
{
    "Total":1842,
    "Periods":60,
    "GreaterIntensity":3,
    "GreaterIntensityDays":
    [
        {
            "Wheather":"Rainy",
            "day":2808
        },
        {
            "Wheather":"Rainy",
            "day":2988
        },
        {
            "Wheather":"Rainy",
            "day":3168
        }
    ],
    "Days":
    [
        {
            "Wheather":"Rainy",
            "day":22
        },
        ...
    ]
}
```

* Consultar todos los días (Solo traerá los almacenados en memoria)

```url
http://{host}:Port/wheather/all
```

Ejemplo de Respuesta:

```url
http://localhost:8080/wheather/all
```

```json
{
    "Total":3656,
    "Days":[
        {
            "Wheather":"Normal",
            "day":-3
        },
        {
            "Wheather":"Normal",
            "day":-2
        },{
            "Wheather":"Normal",
            "day":-1
        },
        {
            "Wheather":"Normal",
            "day":1
        }
        ...
    ]
}
```

## Nota

Los datos de la distancia al sol de los planetas, los grados por día, el angulo inicial de los planetas con respecto al eje x, si el movimiento es horario o anti-horario y la fecha en que se inicio el calculo, pueden ser modificados mediante variables de entorno antes de ejecutar el sistema:

```sh
    SOLAR_SYSTEM_INITIAL_DATE
    FERENGINAR_DEGREES_PER_DAY
    FERENGINAR_SUN_DISTANCE
    FERENGINAR_INITIAL_DEGREES
    FERENGINAR_CLOCKWISE
    BETAZED_DEGREES_PER_DAY
    BETAZED_SUN_DISTANCE
    BETAZED_INITIAL_DEGREES
    BETAZED_CLOCKWISE
    VULCANO_DEGREES_PER_DAY
    VULCANO_SUN_DISTANCE
    VULCANO_INITIAL_DEGREES
    VULCANO_CLOCKWISE
```