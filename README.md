# MID Xchango

Repo del MID (Modelo de Intercambio / Módulo de proyecto) - Equipo Xchango.

## Ramas

- `main` -> rama principal, estable. Solo se fusiona desde `release` o hotfix.
- `develop` -> rama de integración. Aquí se juntan todas las features.
- `release` -> rama de preparación de entrega / pruebas finales antes de pasar a `main`.
- `feature/yeferson`, `feature/miguel`, `feature/tomas`, `feature/stiven` -> trabajo individual de cada integrante. Nace de `develop` y vuelve a `develop` con Pull Request.

## Flujo de trabajo

1. Actualizar `develop`:
   ```bash
   git checkout develop
   git pull origin develop
   ```
2. Crear / ubicarse en tu rama:
   ```bash
   git checkout -b feature/tu-nombre develop
   # o si ya existe:
   git checkout feature/tu-nombre
   ```
3. Trabajar y commitear:
   ```bash
   git add .
   git commit -m "describe lo que hiciste"
   git push -u origin feature/tu-nombre
   ```
4. Cuando termines, Pull Request de `feature/tu-nombre` -> `develop`.
5. Para entrega: `develop` -> `release` (pruebas) -> `main` (versión final, con tag).

## Convención de commits

- `feat: ...` nueva funcionalidad
- `fix: ...` corrección
- `docs: ...` documentación
- `chore: ...` tareas, config

## Integrantes

- Yeferson
- Miguel
- Tomas
- Stiven
